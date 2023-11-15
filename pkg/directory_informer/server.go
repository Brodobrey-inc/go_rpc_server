package directoryinformer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testGRPC/pkg/api"
	"testGRPC/pkg/cache"
)

type DirectoryInformer struct {
	api.UnimplementedDirectoryInformerServer

	cache *cache.Cache
}

func NewDirectoryInformer(cache *cache.Cache) *DirectoryInformer {
	return &DirectoryInformer{
		cache: cache,
	}
}

func (di *DirectoryInformer) Dir(ctx context.Context, req *api.DirectoryRequest) (*api.DirectoryResponse, error) {
	desc, ok := di.cache.Get(req.GetPath())
	if !ok {
		entries, err := os.ReadDir(req.GetPath())
		if err != nil {
			return nil, err
		}

		desc = api.Description{Path: req.GetPath()}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				log.Println(err)
				continue
			}

			unit := api.DirectoryUnit{
				Type: api.File,
				Name: info.Name(),
				Size: info.Size(),
			}
			if info.IsDir() {
				unit.Type = api.Directory
			}

			desc.Elements = append(desc.Elements, unit)
		}

		di.cache.Set(req.GetPath(), desc)
	}

	bytes, err := json.Marshal(desc)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &api.DirectoryResponse{
		DirectoryInfo: bytes,
	}, nil
}
