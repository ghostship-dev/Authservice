module default {
    type Token {
        required account: Account;
        required variant: str {
            constraint one_of("access_token", "refresh_token");
            default := "access_token";
        }
        required value: str;
        required scope: array<str>;
        required revoked: bool {
            default := false;
        }
        required expires_at: datetime;
        index on (.value);
    }

    type Authcode {
        required account: Account;
        required application: OAuthApplication;
        required redirect_uri: str;
        required requested_scope: array<str>;
        required granted_scope: array<str> {
            default := <array<str>>{};
        }
        required code: str;
        required consented: bool {
            default := false;
        }
        required expires_at: datetime;
        index on (.code)
    }
}