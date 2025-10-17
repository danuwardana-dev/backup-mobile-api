package articleSvc

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/repository/postgres"
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/dto/request"
	"backend-mobile-api/model/dto/response"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum/pkgErr"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type articleService struct {
	userRepository    postgres.UserRepository
	articleRepository postgres.ArticleRepository
	clogger           *helpers.CustomLogger
}

func NewArticleService(userRepository postgres.UserRepository, articleRepository postgres.ArticleRepository, clogger *helpers.CustomLogger) ArticleService {
	return &articleService{
		userRepository:    userRepository,
		articleRepository: articleRepository,
		clogger:           clogger,
	}
}

type ArticleService interface {
	InsertArticleService(ctx context.Context, req *request.NewArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	UpdateArticleService(ctx context.Context, req *request.UpdateArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	DeleteArticleService(ctx context.Context, req *request.DeleteArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	SelectArticleListService(ctx context.Context, req *request.SelectListArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
	InternalSelectArticleListService(ctx context.Context, req *request.SelectListArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse
}

func (svc *articleService) InsertArticleService(ctx context.Context, req *request.NewArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {

	txArticle := svc.articleRepository.Tx(ctx)
	err := svc.articleRepository.InsertArticle(ctx, txArticle, &entity.Article{
		Name:           req.NewArticle.Name,
		Url:            req.NewArticle.Url,
		Category:       req.NewArticle.Category,
		ActiveAfterDay: req.NewArticle.ActiveAfterDay,
	})
	if err != nil {
		txArticle.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
		}
	}

	logData.Success = true
	txArticle.Commit()
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       nil,
	}
}

func (svc *articleService) UpdateArticleService(ctx context.Context, req *request.UpdateArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		article    *entity.Article
		newArticle = entity.Article{}
		err        error
	)
	article, err = svc.articleRepository.SelectByArticleId(ctx, req.Article.Id)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE,
				Message:    pkgErr.RECORD_NOT_FOUND_MSG,
				Error:      "",
				Data:       nil,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "",
			Data:       nil,
		}
	}
	if req.Article.Name != nil {
		newArticle.Name = *req.Article.Name
	}
	if req.Article.Url != nil {
		newArticle.Url = *req.Article.Url
	}
	if req.Article.Category != nil {
		newArticle.Category = req.Article.Category
	}
	if req.Article.ActiveAfterDay != nil {
		newArticle.ActiveAfterDay = req.Article.ActiveAfterDay
	}
	txArticle := svc.articleRepository.Tx(ctx)
	err = svc.articleRepository.UpdateArticle(ctx, txArticle, article, &newArticle)
	if err != nil {
		txArticle.Rollback()
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "",
			Data:       nil,
		}
	}
	txArticle.Commit()
	logData.Success = true
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       nil,
	}

}

func (svc *articleService) DeleteArticleService(ctx context.Context, req *request.DeleteArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		article *entity.Article
		err     error
	)
	article, err = svc.articleRepository.SelectByArticleId(ctx, req.Article.Id)
	if err != nil {
		logData.Error = err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.ARTICLE_RECORD_NOT_FOUND_CODE,
				Message:    pkgErr.RECORD_NOT_FOUND_MSG,
				Error:      "",
				Data:       nil,
			}
		}
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "",
			Data:       nil,
		}
	}
	txArticle := svc.articleRepository.Tx(ctx)
	err = svc.articleRepository.DeleteArticle(ctx, txArticle, article)
	if err != nil {
		logData.Error = err.Error()
		txArticle.Rollback()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      "",
			Data:       nil,
		}
	}

	logData.Success = true
	txArticle.Commit()
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       nil,
	}
}

func (svc articleService) SelectArticleListService(ctx context.Context, req *request.SelectListArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		user     *entity.User
		articles []entity.Article
		total    int64
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errtmp error
		user, errtmp = svc.userRepository.SelectUserByUUID(ctx, *userUUID)
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errtmp
	})
	g.Go(func() error {
		var errtmp error
		articles, errtmp = svc.articleRepository.SelectListArticles(ctx, req.Article)
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errtmp
	})
	g.Go(func() error {
		var errtmp error
		total, errtmp = svc.articleRepository.CountArticle(ctx, &request.CountArticleRequest{
			Category:       req.Article.Category,
			ActiveAfterDay: req.Article.ActiveAfterDay,
		})
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				total = 0
				return nil
			}
		}
		return errtmp
	})
	err := g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
			Data:       nil,
		}
	}
	if user == nil {
		logData.Error = "user not found"
		return &dto.BaseResponse{
			StatusCode: pkgErr.ARTICLE_USER_NOT_FOUNDCODE,
			Message:    pkgErr.USER_NOT_FOUND_MSG,
		}
	}
	logData.UserUUID = user.UUID
	logData.Email = user.Email
	if user.DeviceID != req.DeviceID {
		logData.Error = "device id not match"
		return &dto.BaseResponse{
			StatusCode: pkgErr.ARTICLE_DEFERENCE_DEVICE_CODE,
			Message:    pkgErr.DEFERENCE_DEVICE_MSG,
		}
	}
	var articleResponse []response.ArticleResponse
	for _, article := range articles {
		articleResponse = append(articleResponse, response.ArticleResponse{
			ID:             article.ID,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      &article.UpdatedAt,
			Name:           article.Name,
			Url:            article.Url,
			Category:       article.Category,
			ActiveAfterDay: article.ActiveAfterDay,
		})
	}

	logData.Success = true
	resData := response.SelectArticleResponse{
		Total:    total,
		Limit:    req.Article.Limit,
		Offset:   req.Article.Offset,
		Articles: articleResponse,
	}
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       resData,
	}
}
func (svc *articleService) InternalSelectArticleListService(ctx context.Context, req *request.SelectListArticleRequest, userUUID *string, logData *dto.CustomLoggerRequest) *dto.BaseResponse {
	var (
		articles        []entity.Article
		total           int64
		ArticleResponse []response.ArticleResponse
	)
	g := errgroup.Group{}
	g.Go(func() error {
		var errtmp error
		articles, errtmp = svc.articleRepository.SelectListArticles(ctx, req.Article)
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				return nil
			}
		}
		return errtmp
	})
	g.Go(func() error {
		var errtmp error
		total, errtmp = svc.articleRepository.CountArticle(ctx, &request.CountArticleRequest{
			Category:       req.Article.Category,
			ActiveAfterDay: req.Article.ActiveAfterDay,
		})
		if errtmp != nil {
			if errors.Is(errtmp, gorm.ErrRecordNotFound) {
				total = 0
				return nil
			}
		}
		return errtmp
	})
	err := g.Wait()
	if err != nil {
		logData.Error = err.Error()
		return &dto.BaseResponse{
			StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
			Message:    pkgErr.SERVER_BUSY,
			Error:      err.Error(),
			Data:       nil,
		}
	}
	logData.Success = true
	for _, article := range articles {
		ArticleResponse = append(ArticleResponse, response.ArticleResponse{
			ID:             article.ID,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      &article.UpdatedAt,
			Name:           article.Name,
			Url:            article.Url,
			Category:       article.Category,
			ActiveAfterDay: article.ActiveAfterDay,
		})
	}

	resData := response.SelectArticleResponse{
		Total:    total,
		Limit:    req.Article.Limit,
		Offset:   req.Article.Offset,
		Articles: ArticleResponse,
	}
	return &dto.BaseResponse{
		StatusCode: pkgErr.SUCCESS_CODE,
		Message:    pkgErr.SUCCES_MSG,
		Error:      "",
		Data:       resData,
	}
}
