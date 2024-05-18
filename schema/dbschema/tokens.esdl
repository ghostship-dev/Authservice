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
}