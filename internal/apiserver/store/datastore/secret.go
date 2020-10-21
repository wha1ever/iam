// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datastore

import (
	"gorm.io/gorm"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/fields"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/iam/internal/pkg/util/gormutil"
)

type secrets struct {
	db *gorm.DB
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds.DB}
}

// Create creates a new secret account.
func (s *secrets) Create(secret *v1.Secret, opts metav1.CreateOptions) error {
	return s.db.Create(&secret).Error
}

// Update updates an secret information by the secret identifier.
func (s *secrets) Update(secret *v1.Secret, opts metav1.UpdateOptions) error {
	return s.db.Save(secret).Error
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(username, name string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}

	return s.db.Where("username = ? and name = ?", username, name).Delete(&v1.Secret{}).Error
}

// DeleteCollection batch deletes the secrets.
func (s *secrets) DeleteCollection(username string, names []string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}

	return s.db.Where("username = ? and name in (?)", username, names).Delete(&v1.Secret{}).Error
}

// Get return an secret by the secret identifier.
func (s *secrets) Get(username, name string, opts metav1.GetOptions) (*v1.Secret, error) {
	secret := &v1.Secret{}
	d := s.db.Where("username = ? and name= ?", username, name).First(&secret)

	return secret, d.Error
}

// List return all secrets.
func (s *secrets) List(username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	ret := &v1.SecretList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	if username != "" {
		s.db = s.db.Where("username = ?", username)
	}

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	name, _ := selector.RequiresExactMatch("name")

	d := s.db.Where(" name like ?", "%"+name+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}