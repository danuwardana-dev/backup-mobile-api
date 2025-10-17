package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/entity"
	"context"
	"gorm.io/gorm"
)

type articleRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewArticleRepository(masterDb *gorm.DB, clogger *helpers.CustomLogger) ArticleRepository {
	return &articleRepository{masterDb: masterDb, clogger: clogger}
}

type ArticleRepository interface {
	Tx(ctx context.Context) *gorm.DB
	InsertArticle(ctx context.Context, tx *gorm.DB, news *entity.Article) error
	UpdateArticle(ctx context.Context, tx *gorm.DB, news *entity.Article, updater *entity.Article) error
	DeleteArticle(ctx context.Context, tx *gorm.DB, news *entity.Article) error
	SelectByArticleId(ctx context.Context, newsId uint) (*entity.Article, error)
	SelectListArticles(ctx context.Context, req request.ArticleByList) ([]entity.Article, error)
	CountArticle(ctx context.Context, articleRequest *request.CountArticleRequest) (int64, error)
}

func (r *articleRepository) CountArticle(ctx context.Context, req *request.CountArticleRequest) (int64, error) {
	var (
		clause entity.Article
		count  int64
		err    error
	)
	clause.Category = req.Category
	clause.ActiveAfterDay = req.ActiveAfterDay

	err = r.masterDb.WithContext(ctx).Where(&clause).Find(&[]entity.Article{}).Count(&count).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "CountArticle.masterDb.WithContext(ctx).Where(&clause).Find(&[]entity.Article{}).Count(&count).Error", err)
	}
	return count, err
}

func (r *articleRepository) SelectListArticles(ctx context.Context, req request.ArticleByList) ([]entity.Article, error) {

	var (
		news   []entity.Article
		clause entity.Article
	)

	if req.Category != nil {
		clause.Category = req.Category
	}
	if req.ActiveAfterDay != nil {
		clause.ActiveAfterDay = req.ActiveAfterDay
	}
	err := r.masterDb.WithContext(ctx).Where(&clause).Limit(req.Limit).Offset(req.Offset).Find(&news).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SelectListArticles.masterDb.WithContext(ctx).Find", err)
	}
	return news, err
}

func (r *articleRepository) SelectByArticleId(ctx context.Context, newsId uint) (*entity.Article, error) {
	var news entity.Article
	err := r.masterDb.WithContext(ctx).Where("id = ?", newsId).First(&news).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "SelectByArticleId.masterDb.WithContext(ctx).Where(newsId).First", err)
		return nil, err
	}
	return &news, err
}

func (r *articleRepository) Tx(ctx context.Context) *gorm.DB {
	return r.masterDb.Begin()
}

func (r articleRepository) InsertArticle(ctx context.Context, tx *gorm.DB, news *entity.Article) error {
	err := tx.Create(news).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "InsertArticle.tx.Create", err)
	}
	return err
}

func (r *articleRepository) UpdateArticle(ctx context.Context, tx *gorm.DB, news *entity.Article, updater *entity.Article) error {
	err := tx.Model(news).Updates(updater).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "UpdateArticle.tx.Model(article-svc).Updates", err)
	}
	return err
}

func (r *articleRepository) DeleteArticle(ctx context.Context, tx *gorm.DB, news *entity.Article) error {
	err := tx.Delete(news).Error
	if err != nil {
		r.clogger.ErrorLogger(ctx, "DeleteArticle.tx.Delete", err)
	}
	return err
}
