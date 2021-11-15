import "source-map-support/register";
import { APIGatewayProxyEvent, APIGatewayProxyHandler } from "aws-lambda";
import { errorResponse, successResponse } from "@libs/apiGateway";

export const encryptMessage: APIGatewayProxyHandler = async (event: APIGatewayProxyEvent) => {
  try {
    console.log(`event`, event);
    return successResponse('success');
  } catch (error) {
    console.log("LOG error", error);
    return errorResponse(error);
  }
};
