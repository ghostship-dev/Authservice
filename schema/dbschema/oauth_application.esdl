module default {
    type OAuthApplication {
        required client_id: str;
        required client_secret: str;
        required client_name: str {
            constraint exclusive;
        }
        required client_type: str;
        required redirect_uris: array<str> {
            default := <array<str>>{};
        }
        required grant_types: array<str> {
            default := ["authorization_code"];
        }
        required scope: array<str> {
            default := <array<str>>{};
        }
        required client_owner: Account;
        client_description: str;
        client_homepage_url: str;
        client_logo_url: str;
        client_tos_url: str;
        client_privacy_url: str;
        required client_registration_date: datetime;
        required client_status: str {
            constraint one_of("active", "disabled", "suspended");
            default := "enabled";
        }
        client_rate_limits: json {
            default := <json>{}
        }
        index on (.client_id);
    }
}