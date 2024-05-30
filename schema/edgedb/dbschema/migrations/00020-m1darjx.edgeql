CREATE MIGRATION m1darjxigmqhll6eextmlsu24gv7cntumn6e63h6gmwfcfjj5lmkrq
    ONTO m1lhyel7nzow4yivkt7upm7e4adylk5aqbe5nvrs4zc3vkrymxcufa
{
  ALTER TYPE default::Authcode {
      ALTER PROPERTY required_scope {
          RENAME TO requested_scope;
      };
  };
};
