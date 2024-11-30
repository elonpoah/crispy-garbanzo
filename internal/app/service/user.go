package service

import (
	"context"
	commonReq "crispy-garbanzo/common/request"
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	"crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/utils"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (user *system.SysUser, errCode int) {
	if nil == global.FPG_DB {
		return nil, response.InternalServerError
	}
	err := global.FPG_DB.Where("username = ?", u.Username).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, response.PasswordError
		}
	}
	if err != nil {
		return nil, response.UserNotFound
	}
	return user, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

func (userService *UserService) Register(u system.SysUser) (user *system.SysUser, errCode int) {
	if !errors.Is(global.FPG_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return user, response.UserNameAlready
	}
	// 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	err := global.FPG_DB.Create(&u).Error
	if err != nil {
		return nil, response.InternalServerError
	}
	return &u, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserService) ChangePassword(uid int, Password string, newPassword string) (errCode int) {
	var user system.SysUser
	err := global.FPG_DB.Where("id = ?", uid).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(Password, user.Password); !ok {
			return response.OldPasswordError
		}
	}
	if err != nil {
		return response.InternalServerError
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.FPG_DB.Save(&user).Error
	if err != nil {
		return response.InternalServerError
	}
	return response.SUCCESS

}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfo
//@description: 用户信息
//@param: u *model.SysUser
//@return: userInter *model.SysUser,err error

func (userService *UserService) GetUserInfo(uid int) (userInfo *system.SysUser, errCode int) {
	err := global.FPG_DB.Where("id = ?", uid).First(&userInfo).Error
	if err != nil {
		return userInfo, response.InternalServerError
	}
	return userInfo, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserDepositList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserDepositList(info request.UserDepositRecordReq, uid int) (list *[]system.Deposit, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Deposit{}).Where("uid = ?", uid)
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

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Deposit
//@description: 分页获取数据
//@param: info request.UserDepositRecordReq
//@return: address string, err error

func (userService *UserService) Deposit(req request.UserDepositReq) (address string, errCode int) {
	var record system.WalletAddress
	err := global.FPG_DB.Where("uid = ? AND type = ?", req.Uid, req.Type).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = global.FPG_DB.Where("enable = ? AND status = ? AND type = ?", 1, 0, req.Type).First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return address, response.ChannelCrowded
		} else if err != nil {
			return address, response.InternalServerError
		}
		record.Status = 1
		record.Uid = req.Uid
		err = global.FPG_DB.Save(&record).Error
		if err != nil {
			return address, response.InternalServerError
		}
		deposit := system.Deposit{
			Uid:       req.Uid,
			Username:  req.Username,
			Type:      req.Type,
			Amount:    req.Amount,
			ToAddress: record.Address,
			Status:    0,
		}
		err = global.FPG_DB.Save(&deposit).Error
		if err != nil {
			return address, response.InternalServerError
		}
		return record.Address, response.SUCCESS
	} else {
		return record.Address, response.SUCCESS
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Deposit
//@description: 分页获取数据
//@param: info request.UserWithdrawReq
//@return: address string, err error

func (userService *UserService) Withdraw(req request.UserWithdrawReq) (errCode int) {
	var userInfo system.SysUser
	err := global.FPG_DB.Where("id = ?", req.Uid).First(&userInfo).Error
	if err != nil {
		return response.InvalidUserId
	}
	if userInfo.Balance < req.Amount {
		return response.AvailableBalanceNoEnough
	}
	userInfo.FreezeBalance += req.Amount
	userInfo.Balance -= req.Amount

	withdraw := system.Withdrawal{
		Uid:       req.Uid,
		Username:  req.Username,
		Type:      req.Type,
		Amount:    req.Amount,
		ToAddress: req.Address,
	}
	err = global.FPG_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&userInfo).Error; err != nil {
			return err
		}
		if err := tx.Create(&withdraw).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.InternalServerError
	}
	return response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserWithdrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserWithdrawList(info request.UserWithdrawRecordReq, uid int) (list *[]system.Withdrawal, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Withdrawal{}).Where("uid = ?", uid)
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

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserWithdrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserFreeSpinList(info commonReq.PageInfo, uid int) (list *[]system.InviteDuty, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.InviteDuty{}).Where("uid = ?", uid)
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

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UserMakeDrawReq
//@description: 创建抽奖
//@param: info request.UserMakeDrawReq
//@return: address string, err error

func (userService *UserService) MakeDraw(req request.UserMakeDrawReq) (key string, errCode int) {
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
	record := system.MemberDraw{
		BonusType: req.BonusType,
		Uid:       req.Uid,
		Bonus:     req.Bonus,
		Count:     req.Count,
		Username:  req.Username,
		DrawId:    key,
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

func (userService *UserService) GetDrawById(id string) (info *system.MemberDraw, errCode int) {
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

func (userService *UserService) ClaimDrawById(id string, uid int) (bonus float64, errCode int) {
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
		// 奖励池为空，标记活动结束
		global.FPG_REIDS.Set(ctx, statusKey, 2, 0)
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

func (userService *UserService) GetUserDrawList(info request.UserMakeDrawRecordReq, uid int) (list *[]system.MemberDraw, total int64, errCode int) {
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
