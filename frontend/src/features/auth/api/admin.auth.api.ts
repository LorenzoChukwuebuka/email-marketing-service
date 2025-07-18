import { ResponseT } from "../../../interface/api.interface";
import axiosInstance from "../../../utils/api";
import { AdminLoginValues } from "../interface/admin.auth.interface";

class AdminAuthAPI {
  static  async adminlogin(loginValues: AdminLoginValues):Promise<ResponseT> {
        const response = await axiosInstance.post('/admin/auth/login', loginValues)
        return response.data
    }
}

export default AdminAuthAPI