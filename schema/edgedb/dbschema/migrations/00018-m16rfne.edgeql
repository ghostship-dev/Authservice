CREATE MIGRATION m16rfnep7nxq2f4k6ta7iadogiwygmhzilqiyik5fvipntqy7huptq
    ONTO m1ahluoyejqnej6bizt7a4nv7vhuwxyqvfdickcodorexgghpapksq
{
  ALTER TYPE default::OAuthApplication {
      ALTER PROPERTY client_name {
          CREATE CONSTRAINT std::exclusive;
      };
  };
};
