package models

import (
	"fmt"

	"github.com/Korpenter/club/internal/utils"
)

type Profit struct {
	Table *Table
	Sum   int
}

func (p *Profit) String() string {
	return fmt.Sprintf("%d %d %s", p.Table.Id, p.Sum, utils.Format(p.Table.TotalTime))
}
