package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IBeaconsService interface {
	GetBeacons(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error)

	GetBeaconById(ctx context.Context, beaconId string) (data.BeaconV1, error)

	GetBeaconByUdi(ctx context.Context, beaconId string) (data.BeaconV1, error)

	CreateBeacon(ctx context.Context, beacon data.BeaconV1) (data.BeaconV1, error)

	UpdateBeacon(ctx context.Context, beacon data.BeaconV1) (data.BeaconV1, error)

	DeleteBeaconById(ctx context.Context, beaconId string) (data.BeaconV1, error)
}
