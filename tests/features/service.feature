# file: service.feature

# http://service:8085/

# set up limits for ip, password and login - 3 attempts

Feature: HTTP API for anti bruteforce service

    Scenario: Adding ip in white list
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/whitelist/add" 1 times
        Then The response code should be 201
        When I send "POST" request to "http://ab_service:8085/check" 4 times
        Then I receive response - "ok=true"
        And Teardown test for service

    Scenario: Adding ip in black list
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/blacklist/add" 1 times
        Then The response code should be 201
        When I send "POST" request to "http://ab_service:8085/check" 1 times
        Then I receive response - "ok=false"
        And Teardown test for service

    Scenario: IP limit exceeded
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/check" 4 times with the same ip
        Then I receive response - "ok=false"
        And Teardown test for service

    Scenario: Login limit exceeded
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/check" 4 times with the same login
        Then I receive response - "ok=false"
        And Teardown test for service

    Scenario: Password limit exceeded
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/check" 4 times with the same password
        Then I receive response - "ok=false"
        And Teardown test for service

    Scenario: Clear Buckets
        Given Setup test for service
        When I send "POST" request to "http://ab_service:8085/check" 4 times with the same ip
        And I receive response - "ok=false"
        Then I send "POST" request to "http://ab_service:8085/clear" 1 times
        And I send "POST" request to "http://ab_service:8085/check" 1 times
        And I receive response - "ok=true"
        And Teardown test for service
