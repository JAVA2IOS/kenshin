package models

import models "kenshin/models/Excel"

type JDRow struct {
	ErpFile   models.SheetRow // erp表
	CosFile   models.SheetRow // 成本表
	MoneyFile models.SheetRow // 快车表、消耗金额表
}
