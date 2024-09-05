package service

import (
	"crispy-garbanzo/common/request"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	systemReq "crispy-garbanzo/internal/app/models/request"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetHomeRecommand
//@description: 首页推荐
//@return: result map[string][]system.ActivitySession,err error

func (SessionService *SessionService) GetHomeRecommand() (result *map[string][]system.ActivitySession, err error) {
	// 当前时间戳
	now := time.Now().Unix() * 1000

	var hotSessions []system.ActivitySession
	var hugeBonusSessions []system.ActivitySession
	var highTwRateSessions []system.ActivitySession

	// 按 uids 从多到少排序，取前 10 条
	err = global.FPG_DB.
		Where("status = ? AND open_time > ?", 0, now).
		Order("uids DESC,open_time ASC").
		Limit(10).
		Find(&hotSessions).Error
	if err != nil {
		return nil, err
	}

	// 按 ActivytyBonus 从大到小排序，取前 10 条
	err = global.FPG_DB.
		Where("status = ? AND open_time > ?", 0, now).
		Order("activyty_bonus DESC,open_time ASC").
		Limit(10).
		Find(&hugeBonusSessions).Error
	if err != nil {
		return nil, err
	}

	err = global.FPG_DB.
		Where("status = ? AND open_time > ?", 0, now).
		Order("activyty_bonus / activyty_limit_count DESC").
		Limit(10).
		Find(&highTwRateSessions).Error
	if err != nil {
		return nil, err
	}

	// 构造返回结果
	result = &map[string][]system.ActivitySession{
		"hot":        hotSessions,
		"hugebonus":  hugeBonusSessions,
		"hightwrate": highTwRateSessions,
	}

	return result, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionById
//@description: 活动场次详情
//@return: result system.ActivitySession,err error

func (SessionService *SessionService) GetSessionById(id int) (session *system.ActivitySession, err error) {
	err = global.FPG_DB.Where("id = ?", id).First(&session).Error
	return session, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionById
//@description: 活动场次详情
//@return: result system.ActivitySession,err error

func (SessionService *SessionService) BuySessionTicket(id int, uid int) (err error) {
	now := time.Now().Unix() * 1000
	var session system.ActivitySession
	err = global.FPG_DB.Where("status = ? AND open_time > ? AND id = ?", 0, now, id).First(&session).Error
	if err != nil {
		return err
	}
	isGot, err := SessionService.CheckSession(id, uid)
	if err != nil {
		return err
	}
	if isGot {
		return errors.New("请须知，只参加一次")
	}
	if session.Uids == session.ActivytyLimitCount {
		return errors.New("当前场次参与人数已满")
	}
	var user system.SysUser
	err = global.FPG_DB.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return err
	}
	if user.Balance < float64(session.ActivytySpend) {
		return errors.New("余额不足,请充值")
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
		Status:             session.Status,
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
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CheckSession
//@description: 活动场次详情
//@return: isGot bool,err error

func (SessionService *SessionService) CheckSession(id int, uid int) (isGot bool, err error) {
	var record system.GameRecord
	err = global.FPG_DB.Where("session_id = ? AND uid = ?", id, uid).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isGot = false
		return isGot, nil
	}
	if record.ID > -1 {
		isGot = true
	} else {
		isGot = false
	}
	return isGot, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSessionList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetSessionList(info systemReq.SessionListReq) (list *[]system.ActivitySession, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 当前时间戳
	now := time.Now().Unix() * 1000
	db := global.FPG_DB.Model(&system.ActivitySession{}).Where("status = ? AND open_time > ?", 0, now)
	if info.Type == 1 {
		db = db.Order("activyty_bonus DESC,open_time ASC")
	}
	if info.Type == 2 {
		db = db.Order("activyty_bonus / activyty_limit_count DESC")
	}
	if info.Type == 3 {
		db = db.Order("uids DESC,open_time ASC")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetGameHistory
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetGameHistory(info request.PageInfo, uid int) (list *[]system.GameRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.GameRecord{}).Where("uid = ?", uid).Order("created_at ASC")
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}
