const LoginTemplate = () => {
  const hello = (e) => {
    e.preventDefault();
  };

  return (
    <>
      <main className="min-h-screen">
        <div className="grid grid-cols-1 md:grid-cols-12 md:h-screen overflow-hidden">
          <div className="col-span-1 shadow-lg bg-white p-4 md:col-span-6 overflow-y-auto border-gray-950"></div>
          <div className="col-span-1 md:col-span-6 overflow-y-auto bg-[rgb(236,245,255)]">
            <div className="bg-white rounded-lg shadow-sm h-screen md:min-h-screen w-full mx-auto mt-8 md:w-[550px]">
              <h1 className="text-center mt-10 mb-4 text-green-500 text-4xl">
                Logo
              </h1>

              <form onSubmit={hello}>
                <div className="w-full mx-auto max-w-sm mb-8">
                  <label className="form-control w-full mt-5">
                    <div className="label">
                      <span className="label-text text-medium font-medium">
                        Student First Name
                      </span>
                    </div>
                    <input
                      type="text"
                      placeholder=""
                      className="p-3 px-4 rounded-btn border border-gray-500 w-full"
                    />
                  </label>
                  <label className="form-control w-full mt-5">
                    <div className="label">
                      <span className="label-text text-medium font-medium">
                        Student Last Name
                      </span>
                    </div>
                    <input
                      type="text"
                      placeholder=""
                      className="p-3 px-4 rounded-btn border border-gray-500 w-full"
                    />
                  </label>
                  <label className="form-control w-full mt-5">
                    <div className="label">
                      <span className="label-text text-medium font-medium">
                        Student Email
                      </span>
                    </div>
                    <input
                      type="email"
                      placeholder=""
                      className="p-3 px-4 rounded-btn border border-gray-500 w-full"
                    />
                  </label>
                </div>
              </form>
            </div>
          </div>
        </div>
      </main>
    </>
  );
};

export default LoginTemplate;
