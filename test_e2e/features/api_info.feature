Feature: Test info controller

Scenario: Server should return ok and current project version
  When I request "/info"
  Then status code is "200"
  And response field "status" is "OK"
  And response field "version" is "0.0.0-unreleased"
