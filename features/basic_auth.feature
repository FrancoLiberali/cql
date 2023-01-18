Feature: Login as superadmin using the basic authentification

  Scenario: Should be a success on valid credentials
    Given I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | admin                     | string |
    Then I expect status code is "200"

  Scenario: Should be an error on invalid credentials
    Given I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | wrongpassword             | string |
    Then I expect status code is "401"
    And I expect response field "err" is "wrong password"
    And I expect response field "msg" is "the provided password is incorrect"
    And I expect response field "status" is "Unauthorized"

  Scenario: Should be a success if we logout after a successful login
    Given I request "/login" with method "POST" with json
      | key      | value                     | type   |
      | email    | admin-no-reply@badaas.com | string |
      | password | admin                     | string |
    Then I expect status code is "200"
    And I expect response field "username" is "admin"
    And I expect response field "email" is "admin-no-reply@badaas.com"

  Scenario: Should be an error if we try to logout without login first
    When I request "/logout"
    Then I expect status code is "401"
    And I expect response field "err" is "Authentification Error"
    And I expect response field "msg" is "not authenticated"
    And I expect response field "status" is "Unauthorized"
