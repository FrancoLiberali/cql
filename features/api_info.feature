Feature: Test info controller

Scenario: Server should return ok and current project version
  When I request "/info"
  Then I expect status code is "200"
  And I expect response field "status" is "OK"
  And I expect response field "version" is "UNRELEASED"
