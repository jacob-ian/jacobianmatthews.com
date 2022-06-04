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
}

export class HttpService {
  public async post<T = unknown>(request: RequestInput): Promise<T> {
    const { url, body, headers } = request;
    const fetchRequest: RequestInit = {
      method: "POST",
    };

    if (body) {
      fetchRequest.body = JSON.stringify(body);
    }

    if (headers) {
      fetchRequest.headers = headers;
    }

    const res = await this._makeRequest(url, fetchRequest);
    return res.json();
  }

  private async _makeRequest(
    url: string,
    request: RequestInit,
  ): Promise<Response> {
    const res = await this._fetch(url, request);
    if (res.ok) {
      return res;
    }
    const body = await res.text();
    const message = this._getErrorMessage(body);

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
}
