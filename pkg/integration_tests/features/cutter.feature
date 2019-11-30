Feature: Cutter service
  As developer of cutter service
  I want to check interaction between cutter service and remote server

  Scenario: Image is found in cache
    When Client make correct requests to image on remote server
    Then Cutter service should put image in to cache

  Scenario: Remote server does not exist
    When Client make request to non exist server
    Then Cutter service should return 503 http code

  Scenario: Image is not found on remote server
    When Client make request to non exist file
    Then Cutter service should return 404 http code

  Scenario: Image has unsupported extension
    When Client make request to non image file
    Then Cutter service should return 422 http code

  Scenario: Remote server return error
    When Remote server return error for correct request
    Then Cutter service should return 500 http code