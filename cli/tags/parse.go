package tags

import (
	"fmt"
	"strings"
)

func Parse(text string) (tags []string, err error) {
	if text = strings.TrimSpace(text); text == "" {
		tags = make([]string, 0)
		return
	}

	var table = "Aa_B_bC_cD_dE_eF_fG_gH_hI_iJ_jK_kL_lM_mN_nO_oP_pQ_qR_rS_sT_tU_uV_vW_wX_xY_yZ_z_0123456789 "
	for _, tag := range strings.Split(text, ",") {
		if strings.Contains(tag, " ") {
			err = fmt.Errorf("tag %s must not contain empty spaces", tag)
			return
		}
		tag = strings.TrimSpace(tag)
		for _, char := range tag {
			if !strings.Contains(table, string(char)) {
				err = fmt.Errorf("character %s is not allowed in tag, only english alphabet, digits from 0 to 9 and _ are allowd", string(char))
				return
			}
		}
		tags = append(tags, tag)
	}
	return
}
