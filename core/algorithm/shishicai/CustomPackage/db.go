package CustomPackage

import (
	"xmn_2/core/model"
	"fmt"
	"strings"
	"strconv"
)

//
func ChangeA()  {
	//获取数据包
	cPackage := new(model.CustomPackage)
	configPackage := cPackage.Query()
	for i := range configPackage {
		var str string
		str = strings.Replace(configPackage[i].Alias, "2a", "3a", -1)
		str = strings.Replace(configPackage[i].Alias, "2A", "3A", -1)
		fmt.Println(str)
		cPackage.Update("UPDATE custom_package SET alias='" + str + "' WHERE id=" + strconv.Itoa(configPackage[i].Id))
	}
}