package model

import (
	"strconv"
	"fmt"
)

type AlarmRecord struct {
	Id 			int
	AlarmId 	int
	Number 		int
	Cycle 		int
	Title 		string
	CpType 		int
	Position	int
	QNumber		int
	ZNumber		int
	HNumber		int
	CreatedAt 	string
}

// 写入数据
func (ar *AlarmRecord) Insert() {
	sql_str := "insert into `alarm_record` ( `alarm_id`, `number`, `cycle`, `title`, `cp_type`, `position`, `q_num`, `z_num`, `h_num`, `created_at`) values ( '"+ strconv.Itoa(ar.AlarmId) +"', '"+ strconv.Itoa(ar.Number) +"', '"+ strconv.Itoa(ar.Cycle) +"', '"+ ar.Title +"', '"+ strconv.Itoa(ar.CpType) +"', '"+strconv.Itoa(ar.Position)+"', '"+ strconv.Itoa(ar.QNumber) +"', '"+ strconv.Itoa(ar.ZNumber) +"', '"+ strconv.Itoa(ar.HNumber) +"', '"+ar.CreatedAt+"');"
	_, err := DB.Exec(sql_str)
	if err != nil {
		fmt.Println("两连号 报警记录写入失败")
		return
	}
}