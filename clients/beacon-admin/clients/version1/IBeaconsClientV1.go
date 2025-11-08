package clients1

import (
	"context"

	data1 "github.com/Shuv1Wolf/subterra-locate/services/beacon-admin/data/version1"
	cdata "github.com/Shuv1Wolf/subterra-locate/services/common/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IBeaconsClientV1 interface {
	GetBeacons(ctx context.Context, reqctx cdata.RequestContextV1, filter *cquery.FilterParams,
		paging *cquery.PagingParams) (*cquery.DataPage[data1.BeaconV1], error)

	GetBeaconById(ctx context.Context, reqctx cdata.RequestContextV1, beaconId string) (*data1.BeaconV1, error)

	GetBeaconByUdi(ctx context.Context, reqctx cdata.RequestContextV1, udi string) (*data1.BeaconV1, error)

	CreateBeacon(ctx context.Context, reqctx cdata.RequestContextV1, beacon data1.BeaconV1) (*data1.BeaconV1, error)

	UpdateBeacon(ctx context.Context, reqctx cdata.RequestContextV1, beacon data1.BeaconV1) (*data1.BeaconV1, error)

	DeleteBeaconById(ctx context.Context, reqctx cdata.RequestContextV1, beaconId string) (*data1.BeaconV1, error)
}
