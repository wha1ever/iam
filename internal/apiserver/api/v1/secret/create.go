// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	"github.com/AlekSi/pointer"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/log"
)

const maxSecretCount = 10

// Create add new secret key pairs to the storage.
func Create(c *gin.Context) {
	log.Info("create secret function called.", log.String("X-Request-Id", requestid.Get(c)))

	var r v1.Secret

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)
		return
	}

	username := c.GetHeader("username")

	sec, err := store.Client().Secrets().List(username, metav1.ListOptions{
		Offset: pointer.ToInt(0),
		Limit:  pointer.ToInt(-1),
	})

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	if sec.TotalCount >= maxSecretCount {
		core.WriteResponse(c, errors.WithCode(code.ErrReachMaxCount, "secret count: %d", sec.TotalCount), nil)
		return
	}

	// must reassign username
	r.Username = username

	if err := store.Client().Secrets().Create(&r, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, r)
}