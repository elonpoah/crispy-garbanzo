package service

import (
	"context"
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	systemReq "crispy-garbanzo/internal/app/models/request"
	systemRes "crispy-garbanzo/internal/app/models/response"
	"crispy-garbanzo/utils"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type SessionService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHomeRecommand
//@description: 首页推荐
//@return: result map[string][]system.ActivitySession,err error

func (SessionService *SessionService) GetHomeRecommand() (result *map[string][]system.ActivitySession, errCode int) {
	// 当前时间戳
	// 当前时间
	currentUnix := time.Now()
	now := currentUnix.Unix()

	var hotSessions []system.ActivitySession
	var hugeBonusSessions []system.ActivitySession
	var highTwRateSessions []system.ActivitySession
	var freeSessions []system.ActivitySession

	// 按 uids 从多到少排序，取前 10 条
	err := global.FPG_DB.
		Where("status = ? AND open_time > ? AND activyty_spend > ?", 1, now, 0).
		Order("uids DESC,open_time ASC").
		Limit(10).
		Find(&hotSessions).Error
	if err != nil {
		return nil, response.InternalServerError
	}

	// 按 ActivytyBonus 从大到小排序，取前 10 条
	err = global.FPG_DB.
		Where("status = ? AND open_time > ? AND activyty_spend > ?", 1, now, 0).
		Order("activyty_bonus DESC,open_time ASC").
		Limit(10).
		Find(&hugeBonusSessions).Error
	if err != nil {
		return nil, response.InternalServerError
	}

	err = global.FPG_DB.
		Where("status = ? AND open_time > ? AND activyty_spend > ?", 1, now, 0).
		Order("activyty_limit_count ASC").
		Limit(10).
		Find(&highTwRateSessions).Error
	if err != nil {
		return nil, response.InternalServerError
	}

	endOfDay := time.Date(currentUnix.Year(), currentUnix.Month(), currentUnix.Day(), 23, 59, 59, 0, currentUnix.Location())
	endOfDayUnix := endOfDay.Unix()
	err = global.FPG_DB.
		Where("status = ? AND open_time > ? AND open_time <= ? AND activyty_spend = ?", 1, now, endOfDayUnix, 0).
		Limit(1).
		Find(&freeSessions).Error
	if err != nil {
		return nil, response.InternalServerError
	}

	// 构造返回结果
	result = &map[string][]system.ActivitySession{
		"hot":        hotSessions,
		"hugebonus":  hugeBonusSessions,
		"hightwrate": highTwRateSessions,
		"free":       freeSessions,
	}

	return result, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionById
//@description: 活动场次详情
//@return: result system.ActivitySession,err error

func (SessionService *SessionService) GetSessionById(id uint) (session *system.ActivitySession, errCode int) {
	err := global.FPG_DB.Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, response.InternalServerError
	}
	return session, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: BuySessionTicket
//@description: 活动场次详情
//@return: result system.ActivitySession,err error

func (SessionService *SessionService) BuySessionTicket(id uint, uid int) (errCode int) {
	now := time.Now().Unix()
	var session system.ActivitySession
	err := global.FPG_DB.Where("status = ? AND open_time > ? AND id = ?", 1, now, id).First(&session).Error
	if err != nil {
		return response.ActivityNotFound
	}
	isGot, errCode := SessionService.CheckSession(id, uid)
	if errCode != response.SUCCESS {
		return response.InternalServerError
	}
	if isGot {
		return response.ActivityGetIn
	}
	if session.Uids == session.ActivytyLimitCount {
		return response.ActivityFullIn
	}
	var user system.SysUser
	err = global.FPG_DB.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return response.UserNotFound
	}
	if user.Balance < float64(session.ActivytySpend) {
		return response.BalanceNotEnough
	}
	session.Uids += 1
	user.Balance = user.Balance - float64(session.ActivytySpend)
	gameRecord := system.GameRecord{
		Uid:                user.ID,
		Username:           user.Username,
		SessionID:          session.ID,
		ActivytyName:       session.ActivytyName,
		ActivytyBonus:      session.ActivytyBonus,
		ActivytySpend:      session.ActivytySpend,
		ActivytyLimitCount: session.ActivytyLimitCount,
		OpenTime:           session.OpenTime,
		Uids:               session.Uids,
		Status:             1,
	}
	err = global.FPG_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		if err := tx.Save(&session).Error; err != nil {
			return err
		}
		if err := tx.Create(&gameRecord).Error; err != nil {
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
//@function: CheckSession
//@description: 活动场次详情
//@return: isGot bool,err error

func (SessionService *SessionService) CheckSession(id uint, uid int) (isGot bool, errCode int) {
	var record system.GameRecord
	err := global.FPG_DB.Where("session_id = ? AND uid = ?", id, uid).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isGot = false
		return isGot, response.SUCCESS
	}
	if record.ID > 0 {
		isGot = true
	} else {
		isGot = false
	}
	return isGot, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetSessionList(info systemReq.SessionListReq) (list *[]system.ActivitySession, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 当前时间戳
	now := time.Now().Unix()
	db := global.FPG_DB.Model(&system.ActivitySession{}).Where("status = ? AND open_time > ? AND activyty_spend > ?", 1, now, 0)
	if info.Type == 1 {
		db = db.Order("activyty_bonus DESC,open_time ASC")
	}
	if info.Type == 2 {
		db = db.Order("activyty_bonus / activyty_limit_count DESC")
	}
	if info.Type == 3 {
		db = db.Order("uids DESC,open_time ASC")
	}
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
//@function: GetGameHistory
//@description: 分页获取数据
//@param: info systemReq.GameHistoryReq
//@return: err error, list interface{}, total int64

func (userService *UserService) GetGameHistory(req systemReq.GameHistoryReq, uid int) (list *[]system.GameRecord, total int64, errCode int) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.FPG_DB.Model(&system.GameRecord{}).Where("uid = ?", uid).Order("open_time DESC")
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	err := db.Count(&total).Error
	if err != nil {
		errCode = response.InternalServerError
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		errCode = response.InternalServerError
		return
	}
	return list, total, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserSummary
//@description: 汇总
//@param: info systemReq.GameHistoryReq
//@return: result systemRes.UserSummaryResponse err error

func (userService *UserService) GetUserSummary(uid int) (result systemRes.UserSummaryResponse, errCode int) {
	var FreeCount int64
	var SessionCount int64
	err := global.FPG_DB.Model(&system.GameRecord{}).Where("uid = ? AND status = ?", uid, 1).Count(&SessionCount).Error
	if err != nil {
		errCode = response.InternalServerError
		return
	}
	err = global.FPG_DB.Model(&system.InviteDuty{}).Where("uid = ?", uid).Count(&FreeCount).Error
	if err != nil {
		errCode = response.InternalServerError
		return
	}
	result.SessionCount = SessionCount
	result.FreeCount = FreeCount

	return result, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CheckInviteDuty
//@description: 邀请注册活动
//@return: result systemRes.InviteSessionResponse,err error

func (userService *UserService) CheckInviteDuty(rangeType int, uid int) (result systemRes.InviteSessionResponse, errCode int) {
	startOfTime, endOfTime := utils.GetTimeRange(rangeType)
	var userIds []uint
	err := global.FPG_DB.Model(&system.SysUser{}).Where("pid = ? AND created_at >= ? AND created_at < ?", uid, startOfTime, endOfTime).Pluck("id", &userIds).Error
	if err != nil {
		return result, response.InternalServerError
	}
	result.Registrations = len(userIds)
	err = global.FPG_DB.Model(&system.GameRecord{}).
		Where("uid IN (?) AND created_at >= ? AND created_at < ?", userIds, startOfTime, endOfTime).
		Select("uid, COUNT(*) as participates").
		Group("uid").
		Count(&result.Participates).Error
	if err != nil {
		return result, response.InternalServerError
	}

	return result, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: StartInviteSpin
//@description: 邀请注册活动
//@return: result int,err error

func (userService *UserService) StartInviteSpin(rangeType int, uid int, inviteValues systemRes.InviteSessionResponse) (bonus float64, errCode int) {

	ctx := context.Background()
	var invite systemRes.InviteConfig
	result, err := global.FPG_REIDS.Get(ctx, system.AppInviteSetting).Result()
	if err != nil {
		return 0, response.InternalServerError
	}
	err = json.Unmarshal([]byte(result), &invite)
	if err != nil {
		return 0, response.InternalServerError
	}
	if rangeType == 1 && (inviteValues.Registrations < int(invite.Daily.Count) || inviteValues.Participates < int64(invite.Daily.Participants)) {
		return 0, response.FreeSpinUnavilable
	}
	if rangeType == 2 && (inviteValues.Registrations < int(invite.Week.Count) || inviteValues.Participates < int64(invite.Week.Participants)) {
		return 0, response.FreeSpinUnavilable
	}
	if rangeType == 3 && (inviteValues.Registrations < int(invite.Month.Count) || inviteValues.Participates < int64(invite.Month.Participants)) {
		return 0, response.FreeSpinUnavilable
	}
	startOfTime, endOfTime := utils.GetTimeRange(rangeType)
	var recordLen int64
	err = global.FPG_DB.Model(&system.InviteDuty{}).Where("type = ? AND uid = ? AND created_at >= ? AND created_at < ?", rangeType, uid, startOfTime, endOfTime).Count(&recordLen).Error
	if err != nil {
		errCode = response.InternalServerError
		return
	}
	if recordLen != 0 {
		return 0, response.FreeSpinAlreadyJoin
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var maxBonus float64

	// 根据 rangeType 确定最大值
	if rangeType == 1 {
		maxBonus = invite.Daily.Bonus
	} else if rangeType == 2 {
		maxBonus = invite.Week.Bonus
	} else if rangeType == 3 {
		maxBonus = invite.Month.Bonus
	}

	// 生成浮点数随机值
	bonus = r.Float64() * maxBonus

	record := system.InviteDuty{
		Uid:    uid,
		Type:   rangeType,
		Amount: bonus,
	}
	err = global.FPG_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		var user system.SysUser
		if err := tx.Where("id = ?", uid).First(&user).Error; err != nil {
			return err
		}
		user.Balance += float64(bonus)
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, response.InternalServerError
	}
	return bonus, response.SUCCESS
}
