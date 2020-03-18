#!/bin/bash

# Globals
access_token=${BD_TOKEN}
web_url="https://blackduck.eng.netapp.com"
version_id="id:9a92dcad-4f31-4d11-ac4c-36e1351a2649"
#component_name=
#component_version=

function authenticate()
{
        response=$(curl --insecure -X POST --header "Content-Type:application/json" --header "Authorization: token $access_token" "$web_url/api/tokens/authenticate")
        bearer_token=$(echo "${response}" | jq --raw-output '.bearerToken')
        if [ ! -z bearer_token ]; then
          echo "Got bearer token from Black Duck hub."
        fi
}

function findAndIgnoreComponent()
{
        # Find all components in a BOM
        command="curl --insecure -X GET --header \"Content-Type:application/json\" --header \"Authorization: bearer $bearer_token\" \"$web_url/api/v1/releases/$version_id/component-bom-entries\""
        echo $command

        components=$(curl --insecure -X GET --header "Content-Type:application/json" --header "Authorization: bearer $bearer_token" "$web_url/api/v1/releases/$version_id/component-bom-entries")
		echo $components

        componentsPretty=$(echo "$components" | jq .)

        echo $componentsPretty
        
        # Find the component using name and version
        #componentToIgnore=$(echo "$components" | jq -c --arg component "$component_name" --arg version "$component_version" '.items[] | select(.projectName==$component and .releaseVersion==$version)')
        
        # Check if we found the component we are looking for
        if [ -z "$componentToIgnore" ]
        then
                # Failed log and exit
                #echo "FAILED: Component $component_name version $component_version not found"
                echo "FAILED: Component version not found"
                exit
        else
                # Found component set the ignore flag to true
                #result=$(echo "$componentToIgnore" | jq -c '.ignored |= true' | jq '[.]')

                # Send the updated component information  
                #out=$(curl --insecure -X PUT --header "Content-Type:application/json" --header "Authorization: bearer $bearer_token" -d "$result" "$web_url/api/v1/releases/$version_id/component-bom-entries")
                out=1
                if [ "$out" -eq "1" ]
                then
                        # Success 
                        #echo "SUCCESS: Ignored component $component_name version $component_version"
                        #echo "Component $component_name version $component_version"
                        echo "Component version not found"
                fi
        fi
}

# Example for me to remember
# https://blackduck.eng.netapp.com/ui/versions/id:9a92dcad-4f31-4d11-ac4c-36e1351a2649/view:bom?sortField=projectName&ascending=true&offset=0&inUseOnly=true

################ MAIN ##################

# Get inputs
authenticate
findAndIgnoreComponent

