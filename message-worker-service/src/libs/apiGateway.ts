import StatusCode from '@enums/StatusCode';

const successResponse = (response: string, statusCode: StatusCode = StatusCode.SUCCESS) => ({
  statusCode,
  body: JSON.stringify(response),
})

const errorResponse = (
    message: string,
    statusCode: StatusCode = StatusCode.SERVER_ERROR,
) => ({
  statusCode,
  body: JSON.stringify(message),
})

export { errorResponse, successResponse }
