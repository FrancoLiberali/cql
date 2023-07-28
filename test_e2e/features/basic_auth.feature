Feature: Login as superadmin using the basic authentication

  Scenario: Should be a success on valid credentials
    When I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | admin                     | string |
    Then status code is "200"
    And response field "username" is "admin"
    And response field "email" is "admin-no-reply@badaas.com"

  Scenario: Should be an error on invalid credentials
    When I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | wrongpassword             | string |
    Then status code is "401"
    And response field "err" is "wrong password"
    And response field "msg" is "the provided password is incorrect"
    And response field "status" is "Unauthorized"

  Scenario: Should be a success if we logout after a successful login
    Given I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | admin                     | string |
    When I request "/logout"
    Then status code is "200"

  Scenario: Should be an error if we try to logout without login first
    When I request "/logout"
    Then status code is "401"
    And response field "err" is "Authentication Error"
    And response field "msg" is "not authenticated"
    And response field "status" is "Unauthorized"
