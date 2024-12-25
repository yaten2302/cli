package get

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/supabase/cli/internal/migration/list"
	"github.com/supabase/cli/internal/utils"
)

func Run(ctx context.Context, branchId string, output string) error {
	resp, err := utils.GetSupabase().V1GetABranchConfigWithResponse(ctx, branchId)
	if err != nil {
		return errors.Errorf("failed to retrieve preview branch: %w", err)
	}
	if resp.JSON200 == nil {
		return errors.New("Unexpected error retrieving preview branch: " + string(resp.Body))
	}

	masked := "******"
	if resp.JSON200.DbUser == nil {
		resp.JSON200.DbUser = &masked
	}
	if resp.JSON200.DbPass == nil {
		resp.JSON200.DbPass = &masked
	}
	if resp.JSON200.JwtSecret == nil {
		resp.JSON200.JwtSecret = &masked
	}

	table := "|HOST|PORT|USER|PASSWORD|JWT SECRET|POSTGRES VERSION|STATUS|"
	if output == "env" {
		table += "|POSTGRES_USER_ENV|"
	}

	table += "\n|-|-|-|-|-|-|-|"
	if output == "env" {
		table += "-|"
	}

	table += "\n"

	row := fmt.Sprintf(
		"|`%s`|`%d`|`%s`|`%s`|`%s`|`%s`|`%s`|",
		resp.JSON200.DbHost,
		resp.JSON200.DbPort,
		*resp.JSON200.DbUser,
		*resp.JSON200.DbPass,
		*resp.JSON200.JwtSecret,
		resp.JSON200.PostgresVersion,
		resp.JSON200.Status,
	)
	if output == "env" {
		row += fmt.Sprintf("`%s`|", "POSTGRES_ENV_URL goes here")
	}
	table += row + "\n"

	return list.RenderTable(table)
}
