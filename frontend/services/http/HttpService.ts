import { BadRequestException } from "./BadRequestException";
import { ForbiddenException } from "./ForbiddenException";
import { HttpException, HttpStatus } from "./HttpException";
import { InternalErrorException } from "./InternalErrorException";
import { NotFoundException } from "./NotFoundException";
import { UnauthenticatedException } from "./UnauthenticatedException";

interface RequestInput {
  url: string;
  body?: Record<string, any>;
  headers?: Headers;
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
}

type PostRequestInput = Omit<RequestInput, "method">;
type GetRequestInput = Pick<RequestInput, "url" | "headers"> & {
  headers?: Headers;
};

export class HttpService {
  public async post<T = unknown>(request: PostRequestInput): Promise<T> {
    return this._makeRequest<T>({ ...request, method: "POST" });
  }

  private async _makeRequest<T>(request: RequestInput): Promise<T> {
    const { url, body, headers, method } = request;
    const fetchRequest: RequestInit = {
      method,
    };

    if (body) {
      fetchRequest.body = JSON.stringify(body);
    }

    fetchRequest.headers = headers || new Headers();
    fetchRequest.headers.set("Content-Type", "application/json");

    const res = await this._fetch(url, fetchRequest);
    if (res.ok) {
      return res.json();
    }
    const resBody = await res.text();
    const message = this._getErrorMessage(resBody);

    const exceptionMap: {
      [status in HttpStatus]: new (message: string) => HttpException;
    } = {
      400: BadRequestException,
      401: UnauthenticatedException,
      403: ForbiddenException,
      404: NotFoundException,
      500: InternalErrorException,
    };

    const exception = exceptionMap[res.status as HttpStatus];
    throw exception
      ? new exception(message)
      : new HttpException(res.status, message);
  }

  private async _fetch(url: string, request: RequestInit): Promise<Response> {
    try {
      return await fetch(url, request);
    } catch (error) {
      throw new Error("An error occurred");
    }
  }

  private _getErrorMessage(body: string): string {
    try {
      const json = JSON.parse(body);
      return json.message || body;
    } catch {
      return body;
    }
  }

  public async get<T = unknown>(request: GetRequestInput): Promise<T> {
    return this._makeRequest<T>({ ...request, method: "GET" });
  }
}
