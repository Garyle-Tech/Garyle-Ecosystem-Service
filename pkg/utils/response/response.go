package response

type Meta struct {
	Message    string      `json:"message"`
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Meta: Meta{
			Message: message,
			Status:  "success",
			Code:    200,
		},
		Data: data,
	}
}

func NewSuccessResponseWithPagination(data interface{}, message string, pagination Pagination) Response {
	return Response{
		Meta: Meta{
			Message:    message,
			Status:     "success",
			Code:       200,
			Pagination: &pagination,
		},
		Data: data,
	}
}

func NewErrorResponse(code int, message string) Response {
	return Response{
		Meta: Meta{
			Message: message,
			Status:  "error",
			Code:    code,
		},
		Data: nil,
	}
}

func ClientError(code int, message string) Response {
	return NewErrorResponse(code, message)
}

func BadRequestError(message string) Response {
	return NewErrorResponse(400, message)
}

func NotFoundError(message string) Response {
	if message == "" {
		message = "Resource not found"
	}
	return NewErrorResponse(404, message)
}

func UnauthorizedError(message string) Response {
	if message == "" {
		message = "Unauthorized access"
	}
	return NewErrorResponse(401, message)
}

func ForbiddenError(message string) Response {
	if message == "" {
		message = "Forbidden access"
	}
	return NewErrorResponse(403, message)
}

func ServerError(message string) Response {
	if message == "" {
		message = "Internal server error"
	}
	return NewErrorResponse(500, message)
}
