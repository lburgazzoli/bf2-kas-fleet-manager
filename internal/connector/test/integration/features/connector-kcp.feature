Feature: connector kcp

  Background:
    Given the path prefix is "/api/connector_mgmt"

    # User for eval organization id 13640210 configured in internal/connector/test/integration/feature_test.go:27
    Given an org admin user named "Gru" in organization "13640210"
    Given I store userid for "Gru" as ${gru_user_id}

    # eval users used in public API
    Given a user named "Stuart" in organization "13640221"
    Given I store userid for "Stuart" as ${stuart_user_id}
    Given a user named "Kevin" in organization "13640222"
    Given I store userid for "Kevin" as ${kevin_user_id}
    Given a user named "Carl" in organization "13640223"
    Given I store userid for "Carl" as ${carl_user_id}

    # org admin user used in admin API
    Given an org admin user named "Dr. Nefario" in organization "13640211"
    Given I store userid for "Dr. Nefario" as ${drnefario_user_id}
    Given an admin user named "Ricky Bobby" with roles "connector-fleet-manager-admin-full"

    # eval users used in admin API
    Given a user named "Dave" in organization "13640224"
    Given I store userid for "Dave" as ${dave_user_id}
    Given a user named "Phil" in organization "13640225"
    Given I store userid for "Phil" as ${phil_user_id}
    Given a user named "Tim" in organization "13640226"
    Given I store userid for "Tim" as ${tim_user_id}

    # users in organization 13640230
    Given an org admin user named "Dusty" in organization "13640230"
    Given I store userid for "Dusty" as ${dusty_user_id}
    Given a user named "Lucky" in organization "13640230"
    Given I store userid for "Lucky" as ${lucky_user_id}
    Given a user named "Ned" in organization "13640230"
    Given I store userid for "Ned" as ${ned_user_id}

    # users in organization 13640231
    Given an org admin user named "El Guapo" in organization "13640231"
    Given I store userid for "El Guapo" as ${guapo_user_id}

    # agent user
    Given a user named "Gru_shard"

  Scenario: Create namespaces in cluster for organization 13640230
    Given I am logged in as "Dusty"
    When I POST path "/v1/kafka_connector_clusters" with json body:
     """
     {
      "name": "Dusty's Cluster"
     }
     """
    Then the response code should be 202
    And the ".status.state" selection from the response should match "disconnected"

    Given I store the ".id" selection from the response as ${connector_cluster_id}
    When I GET path "/v1/kafka_connector_clusters/${connector_cluster_id}/addon_parameters"
    Then the response code should be 200
    And get and store access token using the addon parameter response as ${shard_token} and clientID as ${clientID}
    And I remember keycloak client for cleanup with clientID: ${clientID}

   # verify user can search for cluster using state and ilike
    When I GET path "/v1/kafka_connector_clusters/?search=id=${connector_cluster_id}+and+name+ilike+dusty%25+and+state+ilike+Disconnected&orderBy=id%2Cstate+desc%2Ccreated_at+asc"
    Then the response code should be 200
    And the ".items[0].id" selection from the response should match "${connector_cluster_id}"

   # There should be default namespace at first for organization 13640230
    Given I am logged in as "Lucky"
    When I GET path "/v1/kafka_connector_namespaces"
    Then the response code should be 200

   # Create an organisation namespace
    Given I am logged in as "Dusty"
    When I POST path "/v1/kafka_connector_namespaces/" with json body:
    """
    {
      "name": "shared_namespace",
      "cluster_id": "${connector_cluster_id}",
      "kind": "organisation",
      "annotations": { "connector_mgmt.bf2.org/profile": "default-profile" }
    }
    """
    Then the response code should be 201
    And I store the ".id" selection from the response as ${org_namespace_id}

    # Create a user namespace
    Given I am logged in as "Lucky"
    When I PUT path "/v1/kafka_connector_namespaces/" with json body:
    """
    {
      "name": "Lucky_namespace",
      "cluster_id": "${connector_cluster_id}",
      "kind": "user",
      "annotations": { "connector_mgmt.bf2.org/profile": "default-profile" }
    }
    """
    Then the response code should be 201

    # Re-Create a user namespace
    Given I am logged in as "Lucky"
    When I PUT path "/v1/kafka_connector_namespaces/" with json body:
    """
    {
      "name": "Lucky_namespace",
      "cluster_id": "${connector_cluster_id}",
      "kind": "user",
      "annotations": { "connector_mgmt.bf2.org/profile": "default-profile" }
    }
    """
    Then the response code should be 201
    And I store the ".id" selection from the response as ${lucky_namespace_id}

   # All organization members MUST be able to see the org tenant namespaces
    Given I am logged in as "Ned"
    When I GET path "/v1/kafka_connector_namespaces/?orderBy=name"
    Then the response code should be 200
    And the ".size" selection from the response should match "2"


    Given I am logged in as "Dusty"
    When I POST path "/v1/kafka_connectors?async=true" with json body:
      """
      {
        "kind": "Connector",
        "name": "example_1",
        "namespace_id": "${org_namespace_id}",
        "connector_type_id": "aws-sqs-source-v1alpha1",
        "kafka": {
          "id":"mykafka",
          "url": "kafka.hostname"
        },
        "schema_registry": {
          "id": "myregistry",
          "url": "registry.hostname"
        },
        "service_account": {
          "client_id": "myclient",
          "client_secret": "test"
        },
        "connector": {
            "aws_queue_name_or_arn": "test",
            "aws_access_key": "test",
            "aws_secret_key": "test",
            "aws_region": "east",
            "kafka_topic": "test_post"
        }
      }
      """
    Then the response code should be 202
    And the ".status.state" selection from the response should match "assigning"
    And the ".connector.kafka_topic" selection from the response should match "test_post"

    Given I am logged in as "Dusty"
    When I PUT path "/v1/kafka_connectors?async=true" with json body:
      """
      {
        "kind": "Connector",
        "name": "example_1",
        "namespace_id": "${org_namespace_id}",
        "connector_type_id": "aws-sqs-source-v1alpha1",
        "kafka": {
          "id":"mykafka",
          "url": "kafka.hostname"
        },
        "schema_registry": {
          "id": "myregistry",
          "url": "registry.hostname"
        },
        "service_account": {
          "client_id": "myclient",
          "client_secret": "test"
        },
        "connector": {
            "aws_queue_name_or_arn": "test",
            "aws_access_key": "test",
            "aws_secret_key": "test",
            "aws_region": "east",
            "kafka_topic": "test_put"
        }
      }
      """
    Then the response code should be 202
    And the ".status.state" selection from the response should match "assigning"
    And the ".connector.kafka_topic" selection from the response should match "test_put"


