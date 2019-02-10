package assetservice

import (
	"github.com/FuturICT2/fin4-core/server/datatype"
	"github.com/kjda/exchange/server/apperrors"
	"github.com/kjda/exchange/server/helpers"
)

// FindAll finds all assets
func (db *Service) FindAll(user *datatype.User) ([]datatype.Asset, error) {
	result := []datatype.Asset{}
	rows, err := db.Query(`
		SELECT
			asset.id,
			asset.name,
			asset.symbol,
			asset.description,
			asset.supply,
			asset.creatorId,
			user.usernname,
			asset.minersCounter,
			asset.favoritesCounter,
			asset.ethereumAddress,
			asset.ethereumTransactionAddress,
			asset.createdAt,
			IF(favorites.blockId, TRUE, FALSE),
		FROM asset asset
		LEFT JOIN user user ON asset.creatorId = user.id
		LEFT JOIN asset_favorites favorites ON asset.id=favorites.assetId AND favorites.userId=?
		ORDER BY id DESC`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c datatype.Asset
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Symbol,
			&c.Description,
			&c.Supply,
			&c.CreatorID,
			&c.CreatorName,
			&c.MinersCounter,
			&c.FavoritesCounter,
			&c.EthereumAddress,
			&c.EthereumTransactionAddress,
			&c.CreatedAt,
			&c.CreatedAtHuman,
			&c.DidUserLike,
		)
		if err != nil {
			apperrors.Critical("assetservice:find-all:1", err)
			return nil, err
		}
		c.CreatedAtHuman = helpers.DateToHuman(c.CreatedAt)
		result = append(result, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}