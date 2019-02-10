package assethandlers

import (
	"net/http"

	"github.com/FuturICT2/fin4-core/server/auth"
	"github.com/FuturICT2/fin4-core/server/datatype"
	"github.com/FuturICT2/fin4-core/server/dbservice"
	"github.com/FuturICT2/fin4-core/server/decimaldt"
	"github.com/gin-gonic/gin"
)

//GetAsset handle get asset route
func GetAsset(sc datatype.ServiceContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := auth.MustGetUser(c)
		//TODO this need to be imporoved, should be GET ASSET, should use tm.FindAsset only
		assetID, isOk := dbservice.StringToID(c.Params.ByName("assetId"))
		if !isOk {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		asset, err := sc.AssetService.FindByID(assetID)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		balance, reserved, err := sc.UserService.FindUserBalance(
			user.ID,
			assetID,
		)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		blocks, err := sc.AssetService.GetAssetBlocks(assetID)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		// return unverified blocks for the token oracle and block creator only
		isOracle := sc.AssetService.IsOracle(user.ID, asset.ID)
		filteredBlocks := []datatype.Block{}
		for _, item := range blocks {
			if item.Status == 1 || isOracle || item.UserID == user.ID {
				filteredBlocks = append(filteredBlocks, item)
			}
		}
		miners, err := sc.AssetService.GetAssetMiners(assetID)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, struct {
			ID                  datatype.ID
			Name                string
			Symbol              string
			CreatorID           datatype.ID
			CreatorName         string
			Description         string
			TotalSupply         int64
			UserBalance         decimaldt.Decimal
			UserReservedBalance decimaldt.Decimal
			Blocks              []datatype.Block
			Miners              []datatype.Miner
		}{
			asset.ID,
			asset.Name,
			asset.Symbol,
			asset.CreatorID,
			asset.CreatorName,
			asset.Description,
			asset.Supply,
			balance,
			reserved,
			filteredBlocks,
			miners,
		})
	}
}
