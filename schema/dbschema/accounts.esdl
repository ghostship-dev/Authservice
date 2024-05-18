module default {
    type Account {
        required username: str;
        required email: str {
            constraint exclusive;
            constraint regexp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
        }
        avatar_uri: str;
        status: str {
            constraint one_of("created", "active", "suspended", "deleted");
            default := "created";
        }
        status_description: str;
        status_changed: datetime;
        created_at: datetime;
        otp_secret: str;
        required otp_state: str {
            constraint one_of("disabled", "enabled", "verifying");
            default := "disabled"
        }
        index on (.email);
    }

    type Password {
        required account: Account;
        required email: str {
            constraint exclusive;
            constraint regexp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
        }
        required password: str {
            constraint min_len_value(8);
        }
        last_used: datetime;
        required failed_attempts: int16 {
            default := 0;
        }
        last_failed_attempt: datetime;
        index on (.email);
    }
}