import { Route, Routes } from "react-router-dom";
import UserDashLayout from "./userDashLayout";
import UserDashPage from "../pages/userDashboard/userHomeDashPage";
import { AccountSettingsTemplate, APISettingsDashTemplate, ContactDashTemplate } from "../templates";
import UserMgtSettingsDashTemplate from "../templates/user/SettingsTemplate/userManagementTemplate";
import GroupContactList from "../templates/user/components/contactGroup/groupcontactlistComponent";
import BillingDashTemplate from "../templates/user/BillingTemplate/BillingDashTemplate";
import TemplateBuilderDashComponent from '../templates/user/templateBuilder/templateDashComponent';
import CreateMarketingTemplateDashBoard from "../templates/user/templateBuilder/createMarketingTemplateDashboard";
import CreateTransactionalTemplateDashBoard from "../templates/user/templateBuilder/createTransactionalTemplateDashboard";
import CampaignDashTemplate from "../templates/user/campaign/campaignDashTemplate";
import EditCampaignForm from "../templates/user/components/campaign/editCampaignComponent";
import CampaignReport from "../templates/user/components/campaign/campaignReportComponent";
import AnalyticsTemplateDash from "../templates/user/Analytics/AnalyticsTemplateDash";
import HelpAndSupport from "../templates/supportAndTicket/HelpAndSupport";
import DomainTemplateDash from "../templates/user/SettingsTemplate/domainsTemplatedash";
import DNSAuthenticationRecords from "../templates/user/components/domain/getSingleDomainRecords";
import NotificationList from "../templates/user/notifications/notificationDashTemplate";

const UserDashRoute = () => (
    <Routes>
        <Route path="dash" element={<UserDashLayout />}>
            <Route index element={<UserDashPage />} />
            <Route path="settings/api" element={<APISettingsDashTemplate />} />
            <Route
                path="settings/user-management"
                element={<UserMgtSettingsDashTemplate />}
            />
            <Route
                path="settings/account-management"
                element={<AccountSettingsTemplate />}
            />
            <Route path="contacts" element={<ContactDashTemplate />} />
            <Route path="view-group" element={<GroupContactList />} />
            <Route path="billing" element={<BillingDashTemplate />} />
            <Route path="templates" element={<TemplateBuilderDashComponent />} />
            <Route path="marketing" element={<CreateMarketingTemplateDashBoard />} />
            <Route path="transactional" element={<CreateTransactionalTemplateDashBoard />} />
            <Route path="campaign" element={<CampaignDashTemplate />} />
            <Route path="campaign/edit/:id" element={<EditCampaignForm />} />
            <Route path="campaign/report/:id" element={<CampaignReport />} />
            <Route path="analytics" element={<AnalyticsTemplateDash />} />
            <Route path="support" element={<HelpAndSupport />} />
            <Route path="settings/domain" element={<DomainTemplateDash />} />
            <Route path="settings/domain/records/:id" element={<DNSAuthenticationRecords />} />
            <Route path="notifications" element={<NotificationList/>} />
        </Route>

    </Routes>
);

export { UserDashRoute };
