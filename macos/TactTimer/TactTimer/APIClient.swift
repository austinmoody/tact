import Foundation

class APIClient {

    static let shared = APIClient()

    private static let apiURLKey = "TactAPIURL"
    private static let defaultURL = "http://localhost:2100"

    var baseURL: String {
        get {
            UserDefaults.standard.string(forKey: Self.apiURLKey) ?? Self.defaultURL
        }
        set {
            UserDefaults.standard.set(newValue, forKey: Self.apiURLKey)
        }
    }

    private init() {}

    func createEntry(userInput: String, completion: @escaping (Result<Void, Error>) -> Void) {
        guard let url = URL(string: "\(baseURL)/entries") else {
            completion(.failure(APIError.invalidURL))
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = ["user_input": userInput]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: body)
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse else {
                completion(.failure(APIError.invalidResponse))
                return
            }

            if (200...299).contains(httpResponse.statusCode) {
                completion(.success(()))
            } else {
                let message = data.flatMap { String(data: $0, encoding: .utf8) } ?? "Unknown error"
                completion(.failure(APIError.serverError(statusCode: httpResponse.statusCode, message: message)))
            }
        }.resume()
    }
}

enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case serverError(statusCode: Int, message: String)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid API URL"
        case .invalidResponse:
            return "Invalid response from server"
        case .serverError(let statusCode, let message):
            return "Server error (\(statusCode)): \(message)"
        }
    }
}
