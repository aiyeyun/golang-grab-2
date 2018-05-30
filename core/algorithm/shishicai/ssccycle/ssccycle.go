package ssccycle

import (
	"fmt"
	"xmn_2/core/model"
	"time"
	"log"
	"xmn_2/core/logger"
	"strings"
	"strconv"
	"xmn_2/core/mail"
)

//计算分析结构体
type computing struct {
	packet_a_map   map[string]string
	cpType     int
	cpTypeName string
	code       []string
	position   string
	packet     *model.SscCycle
}


//开始计算
func Calculation()  {
	fmt.Println("时时彩 a连续b 周期")

	//获取开奖号
	cqssc := new(model.Cqssc)
	cqCodes = cqssc.Query("200")

	xjscc := new(model.Xjssc)
	xjCodes = xjscc.Query("200")

	//获取数据包
	cPackage := new(model.SscCycle)
	configPackage := cPackage.Query()

	cq_q3s, cq_z3s, cq_h3s := getFrontCenterAfterCodes(cqsscType)
	xj_q3s, xj_z3s, xj_h3s := getFrontCenterAfterCodes(xjsscType)

	allCodes := &allCpCodes{
		cq_q3s: cq_q3s,
		cq_z3s: cq_z3s,
		cq_h3s: cq_h3s,

		xj_q3s: xj_q3s,
		xj_z3s: xj_z3s,
		xj_h3s: xj_h3s,
	}

	for i := range configPackage {
		go analysis(configPackage[i], allCodes)
	}

}


//解析数据包
func analysis(packet *model.SscCycle, allCodes *allCpCodes)  {
	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("a出现几期的b - 数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("a出现几期的b - 数据包别名: " + packet.Alias + "报警通知非接受时间段内")
		return
	}

	//数据包 a包 解析成map
	slice_dataTxt_package_a := strings.Split(packet.DataTxt, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageA := make(map[string]string)
	for i := range slice_dataTxt_package_a {
		dataTxtMapPackageA[slice_dataTxt_package_a[i]] = slice_dataTxt_package_a[i]
	}


	//重庆前3
	cq_q3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.cq_q3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "前3",
		packet: packet,
	}

	//重庆中3
	cq_z3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.cq_z3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "中3",
		packet: packet,
	}

	//重庆后3
	cq_h3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.cq_h3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "后3",
		packet: packet,
	}

	//新疆前3
	xj_q3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.xj_q3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "前3",
		packet: packet,
	}

	//新疆中3
	xj_z3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.xj_z3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "中3",
		packet: packet,
	}

	//新疆后3
	xj_h3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		code: allCodes.xj_h3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "后3",
		packet: packet,
	}

	go cq_q3.calculate()
	go cq_z3.calculate()
	go cq_h3.calculate()

	go xj_q3.calculate()
	go xj_z3.calculate()
	go xj_h3.calculate()

}

func (md *computing) calculate()  {
	// 周期数
	var cycle_number int = 0

	// a包连续数
	var continuity_number int = 0

	// a连续完 b在规定期数
	var b_number int = 0

	// 重新累计 连续a
	var a_status bool = false

	var log_html string = "<div>腾讯分分彩 a连续周期 包别名: " + md.packet.Alias + " 位置: "+ md.position + "<div>"
	for i := range md.code {
		log_html += "<br/><div>开奖号: " + md.code[i] + "</div>"

		// 检查是否包含a包
		_, in_a := md.packet_a_map[md.code[i]]

		//检查上一期是否 包含A包
		pre_code := ""
		if i != 0 {
			pre_code = md.code[i - 1]
		}

		var pre_in_a bool = false
		//上一期是否包含A包
		if pre_code != "" {
			_, pre_in_a = md.packet_a_map[pre_code]
		}

		log_html += "<div>本期包含a包:"+ strconv.FormatBool(in_a) +"</div>"
		log_html += "<div>上期包含a包:"+ strconv.FormatBool(pre_in_a) +"</div>"

		// 第一期 出现a包 算 1连续
		if i == 0 && in_a == true {
			continuity_number += 1
			log_html += "<div>第一期 就 包含a包 连续数+1 = " + strconv.Itoa(continuity_number) + "</div>"
			continue
		}

		// a包连续到达 阀值
		if continuity_number >= md.packet.Continuity {
			b_number += 1
			log_html += "<div> 等待b包在规定期数内出现 b开始累计 = " + strconv.Itoa(b_number) + "</div>"

			// a包连续完 在b包规定期数内出现了b
			if b_number <= md.packet.Bnumber && in_a == false {
				// 周期清零
				cycle_number = 0
				// 连续清零
				continuity_number = 0
				// b连续期数清零
				b_number = 0
				// 重新计算连续a
				a_status = true
				log_html += "<div> a包连续完 在b包规定期数内出现了b </div>"
				log_html += "<div> 周期值: "+ strconv.Itoa(cycle_number) +" </div>"
				log_html += "<div> a连续值: "+ strconv.Itoa(continuity_number) +" </div>"
				log_html += "<div> b值: "+ strconv.Itoa(b_number) +" </div>"
			}

			// a包连续完 未在b包规定期数内出现了b
			if b_number > md.packet.Bnumber {
				// 周期+1
				cycle_number += 1
				// 连续清零
				continuity_number = 0
				// b连续期数清零
				b_number = 0
				// 重新计算连续a
				a_status = true
				log_html += "<div> a包连续完 未在b包规定期数内出现了b 周期+1 = "+ strconv.Itoa(cycle_number) +" </div>"
			}
			continue
		}

		// 包含a包 且 上一期 包含a包
		if in_a == true && pre_in_a {
			continuity_number += 1
			// 关闭 重新计算连续a
			a_status = false
			log_html += "<div>本期 包含a包 且 上期包含a包 连续数+1 = " + strconv.Itoa(continuity_number) + "</div>"
			continue
		}

		// 本期b包
		if in_a == false {
			// 连续清零
			continuity_number = 0
			// b连续期数清零
			b_number = 0
			// 重新计算连续a
			a_status = true

			log_html += "<div> 本期b包 a连续值 " + strconv.Itoa(continuity_number) + " </div>"
			continue
		}

		// 重新累计连续a
		if i > 0 && in_a == true && a_status == true {
			continuity_number += 1
			log_html += "<div> 重新累计连续a包 " + strconv.Itoa(continuity_number) + " </div>"
			continue
		}
	}

	// 检查是否报警
	if cycle_number >= md.packet.Cycle - 1 && continuity_number == md.packet.Continuity -1 {
		body_html := "<div>腾讯分分彩 a连续b周期 报警 位置: "+ md.position+ " 数据包别名: "+ md.packet.Alias+ " 几A几B: " + strconv.Itoa(md.packet.Continuity) + " A " + strconv.Itoa(md.packet.Bnumber) + " B " + " 当前累计周期数 "+ strconv.Itoa(cycle_number) + " 当前a连续: "+ strconv.Itoa(continuity_number) +"</div>"
		body_html += log_html
		go mail.SendMail("腾讯分分彩 a连续b周期 报警", body_html)
	}
}

//获取 前中后3 开奖号
func getFrontCenterAfterCodes(cpType int) ([]string, []string, []string) {
	q3codes := make([]string, 0)
	z3codes := make([]string, 0)
	h3codes := make([]string, 0)

	//是否属于重复分析
	isRepeat := isRepeat(cpType)
	if !isRepeat {
		//fmt.Println(cpTypeName[cpType], "等待出现最新的号码")
		return q3codes, z3codes, h3codes
	}

	//重庆时时彩
	if cpType == cqsscType {
		for i := range cqCodes {
			q3s := cqCodes[i].One + cqCodes[i].Two + cqCodes[i].Three
			z3s := cqCodes[i].Two + cqCodes[i].Three + cqCodes[i].Four
			h3s := cqCodes[i].Three + cqCodes[i].Four + cqCodes[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	//新疆时时彩
	if cpType == xjsscType {
		for i:= range xjCodes {
			q3s := xjCodes[i].One + xjCodes[i].Two + xjCodes[i].Three
			z3s := xjCodes[i].Two + xjCodes[i].Three + xjCodes[i].Four
			h3s := xjCodes[i].Three + xjCodes[i].Four + xjCodes[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	return q3codes, z3codes, h3codes
}

//是否属于重复分析
func isRepeat(cpType int) bool {

	//数据库最新开奖号
	var newCode string

	//内存中最新开奖号
	var new_code string

	//重庆时时彩
	if cpType == cqsscType {
		//获取本次查询的最新号码
		if len(cqCodes) == 0 {
			return false
		}
		index := len(cqCodes) - 1
		newCode = cqCodes[index].One + cqCodes[index].Two + cqCodes[index].Three + cqCodes[index].Four + cqCodes[index].Five
	}

	// 新疆时时彩
	if cpType == xjsscType {
		//获取本次查询的最新号码
		if len(xjCodes) == 0 {
			return false
		}
		index := len(xjCodes) - 1
		newCode = xjCodes[index].One + xjCodes[index].Two + xjCodes[index].Three + xjCodes[index].Four + xjCodes[index].Five
	}

	/*
	// 天津时时彩
	if cpType == tjsscType {
		//获取本次查询的最新号码
		if len(tjCodes) == 0 {
			return false
		}
		index := len(tjCodes) - 1
		newCode = tjCodes[index].One + tjCodes[index].Two + tjCodes[index].Three + tjCodes[index].Four + tjCodes[index].Five
	}
	*/

	//获取内存中最新的重新开奖号码
	new_code = newsCode.Get(cpType)
	if new_code == newCode {
		return false
	}
	//刷新最新开奖号码
	newsCode.Set(cpType, newCode)
	return true
}