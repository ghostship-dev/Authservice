CREATE MIGRATION m1ksoa72vlh22hosxarsrds6m2sibj5upghn45xq5f4ipcvroc7fmq
    ONTO m1vinisufbbduudedrobtmqlno27epjvrtbkjl6jy4jdvbin4pscda
{
  ALTER TYPE default::Password {
      ALTER PROPERTY failed_attempts {
          SET REQUIRED USING (<std::int16>{});
      };
  };
};
