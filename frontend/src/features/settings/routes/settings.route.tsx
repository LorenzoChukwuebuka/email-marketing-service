import { RouteObject } from "react-router-dom";
import APISettingsDashTemplate from "../templates/apiTemplate";
import AccountSettingsTemplate from "../templates/accountSettingsTemplate";
import DomainTemplateDash from "../templates/domainsTemplate";
import DNSAuthenticationRecords from "../components/domain/getSingleDomainRecords";

const settingsRoute: RouteObject[] = [
    {
        path: "api",
        element: <APISettingsDashTemplate />
    },
    {
        path: "account-management",
        element: <AccountSettingsTemplate />
    },
    {
        path: "domain",
        element: <DomainTemplateDash />
    },
    {
        path: "domain/records/:id",
        element: <DNSAuthenticationRecords />
    }
]


export default settingsRoute