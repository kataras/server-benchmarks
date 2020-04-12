using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using System.Text;

namespace netcore
{
    public class Startup
    {
        public void Configure(IApplicationBuilder app)
        {
            app.UseMiddleware<PlaintextMiddleware>();
        }
    }

    public class PlaintextMiddleware
    {
        private static byte[] BodyBytes = Encoding.UTF8.GetBytes("Index");

        public PlaintextMiddleware(RequestDelegate next)
        {
        }

        public Task Invoke(HttpContext httpContext)
        {
            return httpContext.Response.Body.WriteAsync(BodyBytes, 0, BodyBytes.Length);
        }
    }
}
