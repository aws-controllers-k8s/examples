# ACK resource Health Check in argocd
ArgoCD supports custom health checks written in Lua. 
This document describes how to enable custom health checks for ACK resources in ArgoCD.
## Argo Custom Health Checks
ref: https://argo-cd.readthedocs.io/en/stable/operator-manual/health/#custom-health-checks

add the following to your ArgoCD `argocd-cm` to enable custom health checks for ACK resources:
- `*.services.k8s.aws/*` - for all ACK resources
- `services.k8s.aws/AdoptedResource` - for ACK adopted resource
```yaml
data:
    resource.customizations: |
      services.k8s.aws/AdoptedResource:
        health.lua: |
          hs = {}
          if obj.status ~= nil then
            if obj.status.conditions ~= nil then
              for i, condition in ipairs(obj.status.conditions) do
                if condition.type == "ACK.Adopted" and condition.status == "False" then
                  hs.status = "Degraded"
                  hs.message = condition.message
                  return hs
                end
                if condition.type == "ACK.Adopted" and condition.status == "True" then
                  hs.status = "Healthy"
                  hs.message = condition.message
                  return hs
                end
              end
            end
          end

          hs.status = "Progressing"
          hs.message = "Waiting for Status conditions"
          return hs

      "*.services.k8s.aws/*":
        health.lua.useOpenLibs: true
        health.lua: |
          hs = {}
          if obj.status and obj.status.conditions then
              for i, condition in ipairs(obj.status.conditions) do
                  if condition.type == "ACK.Recoverable" and condition.status == "True" then
                      hs.status = "Degraded"
                      hs.message = condition.message
                      return hs
                  elseif condition.type == "ACK.Terminal" and condition.status == "True" then
                      hs.status = "Degraded"
                      hs.message = condition.message
                      return hs
                  elseif condition.type == "ACK.ResourceSynced" then
                      if condition.status == "True" then
                          hs.status = "Healthy"
                          hs.message = condition.message
                          return hs
                      elseif condition.status == "False" then
                          hs.status = "Progressing"
                          hs.message = condition.reason
                          return hs
                      elseif condition.status == "Unknown" then
                          hs.status = "Degraded"
                          hs.message = condition.reason
                          return hs
                      end
                  end
              end
          end

          hs.status = "Progressing"
          hs.message = "Waiting for Status conditions"
          return hs
```