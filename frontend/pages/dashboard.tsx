import { NextPage } from "next";
import { useAuth } from "../hooks/useAuth";

const Dashboard: NextPage = () => {
  const auth = useAuth({ redirectTo: "/login" });

  if (!auth) {
    return <p>Loading...</p>;
  }

  if (!auth.admin) {
    return <h1>Forbidden</h1>;
  }

  return (
    <>
      <h1>Dashboard</h1>
      {auth}
    </>
  );
};

export default Dashboard;
