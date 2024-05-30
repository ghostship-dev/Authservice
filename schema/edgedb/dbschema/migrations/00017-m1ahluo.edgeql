CREATE MIGRATION m1ahluoyejqnej6bizt7a4nv7vhuwxyqvfdickcodorexgghpapksq
    ONTO m1izmclnysxwtzdusxdnf7o4gsvqx3ps6p66jnckvawg2hfslf6fja
{
  ALTER TYPE default::OAuthApplication {
      ALTER PROPERTY client_rate_limits {
          SET default := (<std::json>{});
      };
  };
};
