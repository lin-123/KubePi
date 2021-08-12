package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(cluster *v1Cluster.Cluster, options common.DBOptions) error
	Update(name string, cluster *v1Cluster.Cluster, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1Cluster.Cluster, error)
	List(options common.DBOptions) ([]v1Cluster.Cluster, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error)
}

func NewService() Service {
	return &cluster{
		DefaultDBService: common.DefaultDBService{},
	}
}

type cluster struct {
	common.DefaultDBService
}

func (c *cluster) Update(name string, cluster *v1Cluster.Cluster, options common.DBOptions) error {
	db := c.GetDB(options)
	r, err := c.Get(name, options)
	if err != nil {
		return err
	}
	cluster.UUID = r.UUID
	cluster.CreateAt = r.CreateAt
	cluster.UpdateAt = time.Now()
	return db.Update(cluster)
}

func (c *cluster) Create(cluster *v1Cluster.Cluster, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster.UUID = uuid.New().String()
	cluster.CreateAt = time.Now()
	cluster.UpdateAt = time.Now()
	return db.Save(cluster)
}

func (c *cluster) Get(name string, options common.DBOptions) (*v1Cluster.Cluster, error) {
	db := c.GetDB(options)
	var cluster v1Cluster.Cluster
	if err := db.One("Name", name, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *cluster) List(options common.DBOptions) ([]v1Cluster.Cluster, error) {
	db := c.GetDB(options)
	var clusters []v1Cluster.Cluster
	if err := db.All(&clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (c *cluster) Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error) {
	db := c.GetDB(options)
	query := db.Select().Limit(size).Skip((num - 1) * size)
	count, err := query.Count(&v1Cluster.Cluster{})
	if err != nil {
		return nil, 0, err
	}
	clusters := make([]v1Cluster.Cluster, 0)
	if err := query.Find(&clusters); err != nil {
		return clusters, 0, err
	}
	return clusters, count, nil
}

func (c *cluster) Delete(name string, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster, err := c.Get(name, options)
	if err != nil {
		return err
	}
	return db.DeleteStruct(cluster)
}
