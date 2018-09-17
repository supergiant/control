import { Injectable } from '@angular/core';
import { UtilService } from '../util/util.service';

@Injectable()
export class AuthService {
  authPath = '/auth';

  constructor(private util: UtilService) { }

  login(data) {
    return this.util.postResponse(this.authPath, data)
      .do((res) => {
        const rawJwt = res.headers.get("authorization");
        this.startSession(rawJwt);
      });
  }

  private startSession(rawJwt) {
    const expiresAt = this.getJwtExpiration(rawJwt);

    localStorage.setItem("authToken", rawJwt);
    localStorage.setItem("expiresAt", expiresAt);
  }


  private getJwtExpiration(rawJwt) {
    const decodedJwt = this.decodeJwt(rawJwt);

    return decodedJwt.expires_at
  }

  private decodeJwt(rawJwt) {
    const payloadUrl = rawJwt.split(".")[1];
    const payload = payloadUrl.replace(/-/g, '+').replace(/_/g, '/');

    return JSON.parse(window.atob(payload))
  }

  logout() {
    localStorage.removeItem("authToken");
    localStorage.removeItem("expiresAt");
  }

  isAuthenticated() {
    const expiresAt = localStorage.getItem("expiresAt");

    return !this.isTokenExpired(expiresAt);
  }

  private isTokenExpired(ms) {
    const currentTime = Date.now().valueOf() / 1000

    return Boolean(ms < currentTime);
  }

  getToken() {
    return localStorage.getItem("authToken");
  }
}
