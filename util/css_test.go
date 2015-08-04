package util

import (
	"fmt"
	"testing"
)

func TestCssCompress1(t *testing.T) {
	css := `
.overall_player{ display:block; }

.phone_download{ width:106px; height:33px;position:absolute; left:290px; top:0px; z-index:999999;}
.phone_download a{ display:block; width:106px; height:33px; background:url(../img/Download.png) no-repeat;}

/* logo */
.logo{ width:280px; height:70px; background:url(../img/logo.png) no-repeat right bottom;position:absolute; left:0px; top:0px; }
.logo a{ display:block;width:280px; height:70px; }

/* Personal  */
.personal{ text-shadow:0 1px 0 #333; height:26px; position:absolute; right:35px; top:0px; z-index:999;}
.personal_left,.personal_right{ width:4px; height:25px; float:left;}
.personal_left{ background:url(../img/bg.png) 0px -36px;}
.personal_right{ background:url(../img/bg.png) -14px -36px;}
.personal_center{ height:25px; float:left; background:url(../img/bg.png); line-height:25px; padding:0 5px;}
.personal_center .set_proset_div{ float:left;}`
	out, err := CssCompress([]byte(css))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
}
