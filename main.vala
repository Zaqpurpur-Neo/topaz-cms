using Gtk;
using WebKit;

public class TopazApp : Gtk.Window {
    private Switch server_switch;
    private Label status_label;
    private WebView web_view;
    private GLib.Subprocess? process = null;

    private string bin_path = "./topaz-server";
    private string server_host = "localhost";
    private uint16 server_port = 7733;
    private string server_url = "http://localhost:7733";

    private int poll_attempts = 0;
    private const int MAX_POLL_ATTEMPTS = 20; // 20 attempts * 200ms = 4 seconds max asumsi

    public TopazApp() {
        this.title = "Topaz Workspace Control";
        this.set_default_size(1100, 750);
        this.window_position = WindowPosition.CENTER;
        this.icon_name = "preferences-system-network";

        web_view = new WebView();
        var settings = web_view.get_settings();
        settings.hardware_acceleration_policy = WebKit.HardwareAccelerationPolicy.ALWAYS;
        settings.enable_webgl = true;
        settings.enable_accelerated_2d_canvas = true;
        settings.enable_developer_extras = true;
        settings.enable_write_console_messages_to_stdout = true;

        var main_box = new Box(Orientation.VERTICAL, 0);


        var header = new HeaderBar();
        header.show_close_button = true;
        header.title = "Topaz Workspace";
        this.set_titlebar(header);

        status_label = new Label("Stopped");
        status_label.margin_end = 10;

        server_switch = new Switch();
        server_switch.valign = Align.CENTER;
        server_switch.notify["active"].connect(on_switch_toggled);

        header.pack_end(server_switch);
        header.pack_end(status_label);

        var scroll = new ScrolledWindow(null, null);
        scroll.add(web_view);
        main_box.pack_start(scroll, true, true, 0);

        load_placeholder("SYSTEM STANDBY", "Deploy local Go Gin server node by toggling the switch above.", false, "stopped");

        this.add(main_box);

        this.destroy.connect(() => {
            stop_server();
            Gtk.main_quit();
        });
    }

    private void load_placeholder(string title, string subtitle, bool show_spinner, string state) {
        string img_path = "file://" + GLib.Environment.get_current_dir() + "/frontend/dist/topaz.webp";

        string badge_text = "OFFLINE";
        string dot_class = "dot-gray";

        if (state == "starting") {
            badge_text = "INITIALIZING";
            dot_class = "dot-green";
        } else if (state == "error") {
            badge_text = "SYSTEM ERROR";
            dot_class = "dot-red";
        }

        string loader_html = show_spinner ? """
            <div class="loader-container">
                <div class="loader-fill"></div>
            </div>
        """ : "";

        // Tactical Game Launcher HUD Template
        string html = """
        <!DOCTYPE html>
        <html>
        <head>
          <meta charset="utf-8">
          <style>
            :root {
              --background-color: #090909;
              --text-color: #fefefe;
              --accent-green: #07b36e;
              --accent-red: #b30749;
              --gray-color-1: #141414;
              --gray-color-2: #282828;
              --text-color-1: #dfdfdf;
              --text-dim: #71717a;
            }

            * { box-sizing: border-box; margin: 0; padding: 0; }

            body {
              font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "JetBrains Mono", monospace;
              background-color: var(--background-color);

              background-image:
                radial-gradient(circle at 50% 50%, rgba(20, 20, 20, 0.8) 0%, var(--background-color) 100%),
                linear-gradient(rgba(40, 40, 40, 0.2) 1px, transparent 1px),
                linear-gradient(90deg, rgba(40, 40, 40, 0.2) 1px, transparent 1px);
              background-size: 100% 100%, 32px 32px, 32px 32px;
              color: var(--text-color);
              display: flex;
              align-items: center;
              justify-content: center;
              height: 100vh;
              overflow: hidden;
              user-select: none;
            }

            /* Launcher Card Frame */
            .launcher-card {
              position: relative;
              width: 520px;
              background-color: var(--gray-color-1);
              border: 1px solid var(--gray-color-2);
              border-radius: 6px;
              padding: 36px 32px 32px 32px;
              box-shadow: 0 25px 50px rgba(0, 0, 0, 0.85), 0 0 40px rgba(7, 179, 110, 0.03);
            }

            /* Cyber Corner Brackets */
            .corner {
              position: absolute;
              width: 10px;
              height: 10px;
              border-color: var(--accent-green);
              border-style: solid;
            }
            .corner-tl { top: -1px; left: -1px; border-width: 2px 0 0 2px; }
            .corner-tr { top: -1px; right: -1px; border-width: 2px 2px 0 0; }
            .corner-bl { bottom: -1px; left: -1px; border-width: 0 0 2px 2px; }
            .corner-br { bottom: -1px; right: -1px; border-width: 0 2px 2px 0; }

            /* Header HUD */
            .launcher-header {
              display: flex;
              align-items: center;
              justify-content: space-between;
              border-bottom: 1px solid var(--gray-color-2);
              padding-bottom: 20px;
              margin-bottom: 24px;
            }

            .brand {
              display: flex;
              align-items: center;
              gap: 14px;
            }

            .logo-box {
              width: 46px;
              height: 46px;
              background-color: var(--background-color);
              border: 1px solid var(--gray-color-2);
              border-radius: 8px;
              display: flex;
              align-items: center;
              justify-content: center;
              padding: 6px;
            }

            .logo-box img {
              max-width: 100%;
              max-height: 100%;
              object-fit: contain;
            }

            .brand-info h1 {
              font-size: 15px;
              font-weight: 700;
              letter-spacing: 0.1em;
              text-transform: uppercase;
              color: var(--text-color);
            }

            .brand-info span {
              font-size: 10px;
              letter-spacing: 0.16em;
              color: var(--text-dim);
              font-family: monospace;
            }

            /* Status Pill Badge */
            .badge {
              display: inline-flex;
              align-items: center;
              gap: 8px;
              padding: 6px 12px;
              border-radius: 4px;
              background-color: var(--background-color);
              border: 1px solid var(--gray-color-2);
              font-size: 10px;
              font-family: monospace;
              font-weight: 700;
              letter-spacing: 0.12em;
            }

            .status-dot {
              width: 7px;
              height: 7px;
              border-radius: 50%;
            }

            .dot-green { background-color: var(--accent-green); box-shadow: 0 0 10px var(--accent-green); animation: pulse 1.8s infinite; }
            .dot-red { background-color: var(--accent-red); box-shadow: 0 0 10px var(--accent-red); }
            .dot-gray { background-color: #444444; }

            @keyframes pulse {
              0% { opacity: 0.4; }
              50% { opacity: 1; }
              100% { opacity: 0.4; }
            }

            /* Main Content Area */
            .main-info {
              margin-bottom: 24px;
              text-align: left;
            }

            .main-info h2 {
              font-size: 18px;
              font-weight: 600;
              color: var(--text-color);
              margin-bottom: 6px;
              letter-spacing: -0.01em;
            }

            .main-info p {
              font-size: 13px;
              color: var(--text-color-1);
              line-height: 1.5;
            }

            /* Tactical Spec Grid */
            .telemetry-grid {
              display: grid;
              grid-template-columns: repeat(3, 1fr);
              gap: 10px;
              background-color: var(--background-color);
              border: 1px solid var(--gray-color-2);
              border-radius: 6px;
              padding: 12px;
              margin-bottom: 20px;
              font-family: monospace;
            }

            .telemetry-item {
              display: flex;
              flex-direction: column;
              gap: 2px;
            }

            .telemetry-label {
              font-size: 9px;
              color: var(--text-dim);
              letter-spacing: 0.12em;
              text-transform: uppercase;
            }

            .telemetry-value {
              font-size: 12px;
              font-weight: 600;
              color: var(--text-color);
            }

            /* Tactical Loader Bar */
            .loader-container {
              width: 100%;
              height: 4px;
              background-color: var(--background-color);
              border-radius: 2px;
              overflow: hidden;
              border: 1px solid var(--gray-color-2);
            }

            .loader-fill {
              height: 100%;
              background-color: var(--accent-green);
              box-shadow: 0 0 12px var(--accent-green);
              animation: loader-slide 1.4s infinite ease-in-out;
            }

            @keyframes loader-slide {
              0% { transform: translateX(-100%); }
              100% { transform: translateX(250%); }
            }
          </style>
        </head>
        <body>
          <div class="launcher-card">
            <div class="corner corner-tl"></div>
            <div class="corner corner-tr"></div>
            <div class="corner corner-bl"></div>
            <div class="corner corner-br"></div>

            <div class="launcher-header">
              <div class="brand">
                <div class="logo-box">
                  <img src="
        """ + img_path + """" alt="Topaz" onerror="this.style.display='none';" />
                </div>
                <div class="brand-info">
                  <h1>TOPAZ WORKSPACE</h1>
                  <span>SYSTEM LAUNCHER v1.0</span>
                </div>
              </div>
              <div class="badge">
                <span class="status-dot """ + dot_class + """"></span>
                """ + badge_text + """
              </div>
            </div>

            <div class="main-info">
              <h2>""" + title + """</h2>
              <p>""" + subtitle + """</p>
            </div>

            <div class="telemetry-grid">
              <div class="telemetry-item">
                <span class="telemetry-label">PORT</span>
                <span class="telemetry-value">7733</span>
              </div>
              <div class="telemetry-item">
                <span class="telemetry-label">HOST</span>
                <span class="telemetry-value">127.0.0.1</span>
              </div>
              <div class="telemetry-item">
                <span class="telemetry-label">ENGINE</span>
                <span class="telemetry-value">GIN / REST</span>
              </div>
            </div>

            """ + loader_html + """
          </div>
        </body>
        </html>
        """;

        web_view.load_html(html, "file://" + GLib.Environment.get_current_dir() + "/");
    }

    private void on_switch_toggled() {
        if (server_switch.active) {
            start_server();
        } else {
            stop_server();
        }
    }

    private void start_server() {
        if (process != null) return;

        try {
            string[] args = { bin_path };
            process = new GLib.Subprocess.newv(args, GLib.SubprocessFlags.NONE);

            status_label.label = "Starting...";
            load_placeholder("DEPLOYING ENGINE", "Initializing Go Gin server process and binding local sockets...", true, "starting");

            poll_attempts = 0;
            GLib.Timeout.add(200, check_server_ready);

        } catch (Error e) {
            server_switch.active = false;
            status_label.label = "Error";
            load_placeholder("DEPLOYMENT FAILED", e.message, false, "error");
            stderr.printf("Failed to start Topaz server: %s\n", e.message);
        }
    }

    private bool check_server_ready() {
        if (process == null || !server_switch.active) {
            return GLib.Source.REMOVE;
        }

        poll_attempts++;

        try {
            var client = new GLib.SocketClient();
            client.timeout = 150;
            var conn = client.connect_to_host(server_host, server_port, null);
            if (conn != null) {
                conn.close();
                status_label.label = "Running";
                web_view.load_uri(server_url);
                return GLib.Source.REMOVE;
            }
        } catch (Error e) {
            // Port not open yet, continue loop
        }

        if (poll_attempts >= MAX_POLL_ATTEMPTS) {
            stop_server();
            status_label.label = "Timeout";
            load_placeholder("CONNECTION TIMEOUT", "Backend process launched, but failed to respond on port 7733.", false, "error");
            return GLib.Source.REMOVE;
        }

        return GLib.Source.CONTINUE;
    }

    private void stop_server() {
        if (process != null) {
            process.force_exit();
            process = null;
            status_label.label = "Stopped";
            load_placeholder("SYSTEM STANDBY", "Deploy local Go Gin server node by toggling the switch above.", false, "stopped");
        }
    }

    public static int main(string[] args) {
        GLib.Environment.set_variable("WEBKIT_DISABLE_COMPOSITING_MODE", "0", true);
        GLib.Environment.set_variable("WEBKIT_DISABLE_DMABUF_RENDERER", "0", true);

        Gtk.init(ref args);
        var app = new TopazApp();
        app.show_all();
        Gtk.main();
        return 0;
    }
}
