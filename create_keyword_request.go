package zensend

import (
  "net/url"
  "strconv"
)


type CreateKeywordRequest struct {
  Shortcode           string
  Keyword             string
  IsSticky            bool
  MoUrl               string
}

func (request *CreateKeywordRequest) toPostParams() url.Values {
  postParams := url.Values{}

  postParams.Set("SHORTCODE", request.Shortcode)
  postParams.Set("KEYWORD", request.Keyword)
  postParams.Set("IS_STICKY", strconv.FormatBool(request.IsSticky))

  if len(request.MoUrl) > 0 {
    postParams.Set("MOURL", request.MoUrl)
  }

  return postParams
}
