import { NextPage } from "next";
import { useAuth } from "../hooks/useAuth";

const Dashboard: NextPage = () => {
  const { user, loading } = useAuth({ redirectTo: "/login" });

  if (loading) {
    return <p>Loading...</p>;
  }

  if (user && !user.admin) {
    return <h1>Forbidden</h1>;
  }

  return (
    <div>
      <h1>Dashboard</h1>
      {JSON.stringify(user)}
    </div>
  );
};

export default Dashboard;
