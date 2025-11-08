package logic

import (
	"context"

	data "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
)

type IBeaconsService interface {
	GetBeacons(ctx context.Context, reqctx cdata.RequestContextV1, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[data.BeaconV1], error)

	GetBeaconById(ctx context.Context, reqctx cdata.RequestContextV1, beaconId string) (data.BeaconV1, error)

	GetBeaconByUdi(ctx context.Context, reqctx cdata.RequestContextV1, beaconId string) (data.BeaconV1, error)

	CreateBeacon(ctx context.Context, reqctx cdata.RequestContextV1, beacon data.BeaconV1) (data.BeaconV1, error)

	UpdateBeacon(ctx context.Context, reqctx cdata.RequestContextV1, beacon data.BeaconV1) (data.BeaconV1, error)

	DeleteBeaconById(ctx context.Context, reqctx cdata.RequestContextV1, beaconId string) (data.BeaconV1, error)
}
