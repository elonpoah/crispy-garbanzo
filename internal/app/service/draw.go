package service

import (
	"context"
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	"crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/utils"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DrawService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UserMakeDrawReq
//@description: 创建抽奖
//@param: info request.UserMakeDrawReq
//@return: address string, err error

func (drawService *DrawService) MakeDraw(req request.UserMakeDrawReq) (key string, errCode int) {
	var userInfo system.SysUser
	err := global.FPG_DB.Where("id = ?", req.Uid).First(&userInfo).Error
	if err != nil {
		return key, response.InvalidUserId
	}

	// 计算总支付金额
	var totalPayAmount float64
	if req.BonusType == 2 { // 平均分配
		totalPayAmount = float64(req.Count) * req.Bonus
	} else { // 随机分配
		totalPayAmount = req.Bonus
	}

	// 检查余额
	if userInfo.Balance < totalPayAmount {
		return key, response.BalanceNotEnough
	}

	// 扣除余额
	userInfo.Balance -= totalPayAmount

	// 生成抽奖唯一键
	key = utils.GenerateUUID12()
	expiresAt := time.Now().Add(60 * time.Minute)
	record := system.MemberDraw{
		BonusType: req.BonusType,
		Uid:       req.Uid,
		Bonus:     req.Bonus,
		Count:     req.Count,
		Username:  req.Username,
		DrawId:    key,
		ExpiresAt: expiresAt,
	}

	// 初始化 Redis 奖励池的 TTL
	const drawTTL = 48 * time.Hour

	// Redis 初始化奖励池函数
	initializeBonusPool := func() error {
		ctx := context.Background()

		// 奖励池键名
		bonusKey := "draw_bonus_pool:" + key
		participantsKey := "draw_participants:" + key
		statusKey := "draw_status:" + key

		// 生成奖励
		var rewards []float64
		if req.BonusType == 2 {
			// 平均分配
			rewards = make([]float64, req.Count)
			for i := 0; i < req.Count; i++ {
				rewards[i] = req.Bonus
			}
		} else {
			// 随机分配
			rewards = utils.GenerateRandomParts(req.Bonus, req.Count)
		}

		// 写入 Redis 列表
		for _, reward := range rewards {
			if err := global.FPG_REIDS.RPush(ctx, bonusKey, reward).Err(); err != nil {
				return err
			}
		}

		// 设置 TTL
		if err := global.FPG_REIDS.Expire(ctx, bonusKey, drawTTL).Err(); err != nil {
			return err
		}

		// 初始化参与者集合和状态
		if err := global.FPG_REIDS.SAdd(ctx, participantsKey, "initialized").Err(); err != nil {
			return err
		}
		if err := global.FPG_REIDS.Expire(ctx, participantsKey, drawTTL).Err(); err != nil {
			return err
		}
		if err := global.FPG_REIDS.Set(ctx, statusKey, 1, drawTTL).Err(); err != nil {
			return err
		}

		return nil
	}

	// 使用事务保存到数据库并初始化 Redis
	err = global.FPG_DB.Transaction(func(tx *gorm.DB) error {
		// 更新用户余额
		if err := tx.Save(&userInfo).Error; err != nil {
			return err
		}

		// 创建抽奖记录
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 初始化 Redis 奖励池
		if err := initializeBonusPool(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return key, response.InternalServerError
	}

	return key, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionById
//@description: 奖励详情
//@return: result system.ActivitySession,err error

func (userService *DrawService) GetDrawById(id string) (info *system.MemberDraw, errCode int) {
	err := global.FPG_DB.Where("draw_id = ?", id).First(&info).Error
	if err != nil {
		return nil, response.InternalServerError
	}
	return info, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ClaimDrawById
//@description: 领取奖励
//@return: result system.ActivitySession,err error

func (drawService *DrawService) ClaimDrawById(id string, uid int) (bonus float64, errCode int) {
	ctx := context.Background()

	// 构造 Redis Key
	bonusKey := fmt.Sprintf("draw_bonus_pool:%s", id)
	participantsKey := fmt.Sprintf("draw_participants:%s", id)
	statusKey := fmt.Sprintf("draw_status:%s", id)
	lockKey := fmt.Sprintf("lock:draw:%s", id)

	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
	success, err := global.FPG_REIDS.SetNX(ctx, lockKey, lockValue, 2*time.Second).Result()
	if err != nil {
		return 0, response.InternalServerError
	}
	if !success {
		return 0, response.ChannelCrowded
	}

	defer func() {
		// 使用 Lua 脚本确保解锁操作的原子性
		unlockScript := `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end
		`
		_, err := global.FPG_REIDS.Eval(ctx, unlockScript, []string{lockKey}, lockValue).Result()
		if err != nil {
			// log.Printf("failed to unlock: %v", err)
			return
		}
	}()

	// 检查活动状态
	status, err := global.FPG_REIDS.Get(ctx, statusKey).Int()
	if err != nil {
		return 0, response.ActivityNotFound // 活动不存在
	}
	if status != 1 {
		return 0, response.ActivityEnded // 活动已结束
	}

	// 检查用户是否已参与
	alreadyParticipated, err := global.FPG_REIDS.SIsMember(ctx, participantsKey, uid).Result()
	if err != nil {
		return 0, response.InternalServerError // Redis 操作失败
	}
	if alreadyParticipated {
		return 0, response.AlreadyParticipated // 已参与过抽奖
	}

	// 获取奖励
	bonus, err = global.FPG_REIDS.LPop(ctx, bonusKey).Float64()
	if err == redis.Nil {
		// 奖励池为空，活动结束
		return 0, response.ActivityEnded
	} else if err != nil {
		return 0, response.InternalServerError // Redis 操作失败
	}

	// 更新用户余额
	err = global.FPG_DB.Transaction(func(tx *gorm.DB) error {
		var userInfo system.SysUser
		// 查询用户信息
		if err := tx.Where("id = ?", uid).First(&userInfo).Error; err != nil {
			return err
		}

		// 更新余额
		userInfo.Balance += bonus
		if err := tx.Save(&userInfo).Error; err != nil {
			return err
		}

		record := system.InviteDuty{
			Uid:    uid,
			Type:   4,
			Amount: bonus,
		}

		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 更新 MemberDraw 表
		var draw system.MemberDraw
		if err := tx.Where("draw_id = ?", id).First(&draw).Error; err != nil {
			return err
		}

		// 累计派奖金额和参与人数
		draw.Distribute += bonus
		draw.Participants++
		// 领取人数=发放个数，设置状态为结束
		if draw.Participants == uint(draw.Count) {
			global.FPG_REIDS.Set(ctx, statusKey, 2, 0)
			draw.Status = 2
		}

		if err := tx.Save(&draw).Error; err != nil {
			return err
		}

		// 记录用户参与
		_, err := global.FPG_REIDS.SAdd(ctx, participantsKey, uid).Result()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_, rollbackErr := global.FPG_REIDS.LPush(ctx, bonusKey, bonus).Result()
		if rollbackErr != nil {
			// log.Printf("failed to rollback bonus to pool: %v", rollbackErr)
			return 0, response.InternalServerError // 数据库事务失败,redis 补偿失败
		}
		return 0, response.InternalServerError // 数据库事务失败
	}

	return bonus, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserDrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (drawService *DrawService) GetUserDrawList(info request.UserMakeDrawRecordReq, uid int) (list *[]system.MemberDraw, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.MemberDraw{}).Where("uid = ?", uid)
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err := db.Count(&total).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	return list, total, response.SUCCESS
}
