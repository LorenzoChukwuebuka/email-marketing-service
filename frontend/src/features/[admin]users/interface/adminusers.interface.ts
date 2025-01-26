import { UserDetails } from '../../auth/interface/auth.interface';

export type AdminUserDetails = UserDetails & { verified_at: string }

export type UserSearchType = "All" | "Verified" | "Unverified"


export type AdminUserStats = {
    total_contacts: number,
    total_campaigns: number,
    total_templates: number,
    total_campaigns_sent: number,
   // total_subscriptions: number,
    total_groups: number
}

